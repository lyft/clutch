---
title: Tailor Your Slack Audit Sink Messages
authors:
  - name: Scarlett Perry
    url: https://github.com/scarlettperry
    avatar: https://user-images.githubusercontent.com/39421794/122237227-2a47ef00-ce8d-11eb-95c9-1b95df5d7706.jpeg
    twitter_username: scarlettp3rry
description: How we extended the Clutch Slack sink to support customized Slack audit messages
image: https://user-images.githubusercontent.com/39421794/122236515-937b3280-ce8c-11eb-9144-67058f4ef78e.png
hide_table_of_contents: false
---

import Image from '@site/src/components/Image';

Imagine being on call for a service and getting paged for high CPU but not knowing why. Traffic hasn't increased and no recent deployments have happened. However, through the Clutch Slack sink, you can see in Slack that a teammate resized your service’s [HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale)/[ASG](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html). Visibility into infrastructure changes is a top level concern for Clutch, which is why the Slack sink is one of the out-of-the-box features.

<!--truncate-->

For each operation performed in Clutch, an audit event is created and can be forward to a channel in your Slack workspace. The Slack sink creates a standard message for all workflows, which answers the questions: what operation was performed, what resources were modified, and who performed the operation.

<Image alt="Clutch Standard Slack Audit Message" src="https://user-images.githubusercontent.com/39421794/122236235-57e06880-ce8c-11eb-8ec7-52963f23c4b9.png" />

The standard message captures key security information, but what if you wanted to add even more information that was specific to a workflow? What if you wanted to customize the information to suit your team's needs? We had the same thoughts and embarked on the next iteration of the Slack sink to support these use cases. Read on to learn more about this recent extension!

## Background

<Image alt="Clutch Audit Architecture" src="https://user-images.githubusercontent.com/39421794/122236515-937b3280-ce8c-11eb-9144-67058f4ef78e.png" variant="centered" />

For each incoming API request into Clutch, an audit event is [created, persisted, and then forwarded](https://clutch.sh/docs/advanced/security-auditing/) to the sinks. The role of a sink is simple: process the audit event appropriately for your needs and forward the event to other systems (Clutch ships with logging and Slack sinks and you can easily add new ones, such as Microsoft Teams). The sinks utilize the recorded information in an audit event to create an output. For example, the Slack sink’s [standard message](https://github.com/lyft/clutch/blob/b56d9a929361fb590144bd8fd25d6b34d0fd5289/backend/service/auditsink/slack/slack.go#L114-L131) is created using the [event's](https://github.com/lyft/clutch/blob/b56d9a929361fb590144bd8fd25d6b34d0fd5289/api/audit/v1/audit.proto#L60-L76) `username`, `service_name`, `method_name`, and `resources` fields.

The first version of the audit middleware (that runs on each request) collected only a subset of information from the incoming API calls. In order to support richer output in the audit sinks (e.g. customized Slack messages), the audit middleware needed to be modified to record the entire API request/response body, which serves as metadata in the audit event.

## The Any

Naturally request and response messages vary across APIs, and the audit middleware doesn’t know the message types in advance (i.e. is this message a [DescribePod](https://github.com/lyft/clutch/blob/a07d897dc69229ad1453fea6dcaf36ad50d2f6a6/api/k8s/v1/k8s.proto#L225-L235) request or a [CreateExperiment](https://github.com/lyft/clutch/blob/a07d897dc69229ad1453fea6dcaf36ad50d2f6a6/api/chaos/experimentation/v1/experimentation.proto#L21-L24) response?). Therefore, we used Protobuf's [Any](https://developers.google.com/protocol-buffers/docs/proto3#any) type to generically serialize the API request/response and save it in the audit event. The Any is composed of two values: the serialized message as bytes and a unique URL for identifying and resolving the message type.

Modifications made to the audit event schema:

```protobuf
message RequestEvent {
  ...
  // The API request saved as metadata for the event.
  RequestMetadata request_metadata = 7;
  // The API response saved as metadata for the event.
  ResponseMetadata response_metadata = 8;
}

message RequestMetadata {
  google.protobuf.Any body = 1;
}

message ResponseMetadata {
  google.protobuf.Any body = 1;
}
```

The event’s Any fields can then be unpacked into its underlying message type by the audit sinks.

Below is a snippet of an audit event with the stored metadata for [DescribePod](https://github.com/lyft/clutch/blob/a07d897dc69229ad1453fea6dcaf36ad50d2f6a6/api/k8s/v1/k8s.proto#L16-L22):

```json
"requestMetadata":{
  "body":{
    "@type":"type.googleapis.com/clutch.k8s.v1.DescribePodRequest",
    "clientset":"kind-clutch-local",
    "cluster":"kind-clutch-local",
    "namespace":"envoy-staging",
    "name":"envoy-58d647f457-5b2k4"
  }
},
"responseMetadata":{
  "body":{
    "@type":"type.googleapis.com/clutch.k8s.v1.DescribePodResponse",
    "pod":{
      "containers":[{"name":"envoy","ready":true,...}],
      "nodeIp":"172.18.0.2",
      "podIp":"10.244.0.5",
      "state":"RUNNING",
      "labels":{"app":"envoy","pod-template-hash":"58d647f457"},
      "podConditions":[{"type":"POD_SCHEDULED","status":"TRUE"},...],
      "status":"Running",
      "startTimeMillis":"1616607169000",
      ...
    }
  }
}
```

## Richer Slack Audits

With new metadata on hand in the audit event, we could support customized Slack messages. Our goals for the feature were to 1) avoid code duplication and 2) make it easy for users to set up for their workflows. Let’s dive into the latter point first.

We [updated](https://github.com/lyft/clutch/blob/c50f64a6a487f7a961a333251afe2cdede3addb7/api/config/service/auditsink/slack/v1/slack.proto#L18-L43) the Slack configuration to allow users to specify a custom message for a `service/method` in the [clutch-config](https://github.com/lyft/clutch/blob/main/backend/clutch-config.yaml). The feature is powered by the Golang [template](https://golang.org/pkg/text/template/) package, so in the freeform message field you’d provide the names of the fields from the workflow’s API request/response. These fields will be replaced by their values when parsed by the Slack sink. We intentionally chose to append the custom message to the standard Slack message to prevent users from removing basic auditing with a less detailed override.

Here’s an example of setting up a message that logs out the min/max size sent in the K8s Resize HPA [request](https://github.com/lyft/clutch/blob/9765a88f042c306b0a9d63fb55a8b2ed552b1c18/api/k8s/v1/k8s.proto#L452-L468):

```yaml
services:
  ...
  - name: clutch.service.db.postgres
  ...
  - name: clutch.service.auditsink.slack
    typed_config:
      “@type”: types.google.com/clutch.config.service.auditsink.slack.v1.SlackConfig
      ...
      filter:
        rules:
          - field: SERVICE
          text: clutch.k8s.v1.K8sAPI
        overrides:
          - full_method: /clutch.k8s.v1.K8sAPI/ResizeHPA
          - message: |
            *Min size*: [[.Request.sizing.min]]
            *Max size*: [[.Request.sizing.max]]
```

A few things to highlight from the example:
  * `.Request.<field name>` is used to retrieve data from the API request.
  * Clutch-specific [templating token](https://github.com/lyft/clutch/blob/0aa1c00b37513900c351be1106cc131498b1aad0/backend/gateway/config.go#L110-L122) is used in lieu of the standard Golang template [Actions](https://golang.org/pkg/text/template/#hdr-Actions) syntax. This is to avoid conflicts during preprocessing of the clutch-config.
  * Slack [markdown](https://api.slack.com/reference/surfaces/formatting#basics) formatting can be added for useful visual highlights to the custom message.
  * Users can provide one custom message per gRPC method.

That minimal configuration is all that’s needed to set up this feature for any Clutch workflow!

<Image alt="Clutch Custom Slack Audit Message" src="https://user-images.githubusercontent.com/39421794/122236354-734b7380-ce8c-11eb-97e4-146222eda962.png" />

## Under the Hood

In the Slack sink, we check if a user has provided a custom Slack message for a given workflow  event and if so, inject the metadata values from the audit event into the custom message and [execute](https://golang.org/pkg/text/template/#Template.Execute) the message template.

As mentioned earlier, one of our goals for the customized Slack message feature was to reduce code duplication. By saving the API request/response body as Any type in the audit event, we can retrieve the metadata values for any custom message with just [this flow](https://github.com/lyft/clutch/blob/9765a88f042c306b0a9d63fb55a8b2ed552b1c18/backend/service/auditsink/slack/slack.go#L174-L202) in the Slack sink. Let’s focus on the important pieces.

```go
reqJSON, err := protojson.Marshal(event.RequestMetadata.Body)
respJSON, err := protojson.Marshal(event.ResponseMetadata.Body)
```

The aim here is to unpack the metadata (Any type) into its underlying message. We looked into using anypb’s [UnmarshalNew](https://pkg.go.dev/google.golang.org/protobuf/types/known/anypb#Any.UnmarshalNew) but ultimately decided upon protojson’s [Marshal](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson#Marshal). The main reason for this choice was that we wanted to use the json format as the basis for the context in order to simplify the metadata structure a bit. For example, protojson flattens [`oneOf`](https://developers.google.com/protocol-buffers/docs/proto3#oneof) fields into its top level field. Since custom messages need to specify what exact data is needed from the request/response, simplifications like these can help make writing a custom message easier.

```go
var requestMetadata map[string]interface{}
var responseMetadata map[string]interface{}

json.Unmarshal(reqJSON, &requestMetadata)
json.Unmarshal(respJSON, &responseMetadata)
```

API request/response messages can be complex types, so we unmarshaled the metadata into a `map[string]interface{}`. This form can be used with many other packages and in particular we wanted to support any of the Golang template [functions](https://golang.org/pkg/text/template/#hdr-Functions) (and make it easy for contributors to enable template helpers like [Sprig](http://masterminds.github.io/sprig/)) so that users can easily retrieve the metadata for their custom Slack audits.

```go
return &auditTemplateData{
	Request:  requestMetadata,
	Response: responseMetadata,

}
```

Lastly, we didn’t want the users to need to know anything about the schema of the audit event metadata, so we saved the metadata in a new simple struct. In this way, the user can chain the information they wanted to retrieve with just `.Request.<field name>` or `.Response.<field name>`.

## Conclusion

The Clutch team values building features as generically as possible so that they can serve many current and future use cases. We didn’t want workflow admins/creators to need to add custom logic to the audit architecture per Clutch workflow to support this feature. Our aim was for it to be intuitive for someone to use without needing to know what’s happening under the hood (like the specifics of the audit flow or Slack sink).

What’s next? In future iterations, we’d like to move the custom messages [alongside API code](https://github.com/lyft/clutch/issues/1519) so that there’s less configuration in the clutch-config and experiment with creating a [standard customized Slack message](https://github.com/lyft/clutch/issues/1520) for workflows. We would also like to add [targeted Slack notifications](https://github.com/lyft/clutch/issues/803), so that you can route your audits to various channels in your Slack workspace.

## Want to get involved?

To learn more about Clutch, contribute, or follow along:
  * Check our [documentation](https://clutch.sh/) and [code](https://github.com/lyft/clutch)
  * Join the [Community](https://clutch.sh/docs/community/)
