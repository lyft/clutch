---
title: Engagement Tracking throuhg Audit Logs
authors:
  - name: Tuong La
    url: https://github.com/lyft-tla
    avatar: https://avatars.githubusercontent.com/u/72420441?v=4
description: Using Clutch's security audit logs to track usage and engagement
image: https://raw.githubusercontent.com/lyft/clutch/main/docs/_website/static/img/docs/audit-architecture-diagram.png
hide_table_of_contents: false
---

import Image from '@site/src/components/Image';

Like most companies, Lyft is data driven. We have hypotheses and use data to support or invalidate them. If you are reading this blog, you probably have a Clutch gateway instance running and have a custom workflow built. If you ever wondered how often those workflows are being used, we will show you how you can easily use Clutch’s built-in Security Auditing as a way to track and report on usage. The audit middleware saves data on each incoming request to a database and any additional sinks of your choosing. With this data, we can track usage of Clutch itself and its integrations. 

Here’s a recap of Clutch’s security auditing architecture ([Security Auditing](https://clutch.sh/docs/advanced/security-auditing)):

<Image alt="Clutch Component Architecture" src="https://raw.githubusercontent.com/lyft/clutch/main/docs/_website/static/img/docs/audit-architecture-diagram.png" width="50%" />

In Lyft’s configuration, we are storing our audit data to a PostgreSQL database in a table called `public.audit_events`. Within `AUDIT_EVENTS`, we want to focus on two columns, `OCCURRED_AT` and `DETAILS`. 

`OCCURRED_AT` stores the date / time of the event.
`DETAILS` stores event information in JSON

If you are not familiar with traversing JSON with Postgres, take a quick peek at JSON Functions and Operators. 

Below is a sample of data that is stored in the `DETAILS` column.

With this SQL statement, we are able to get a list of all users and the actions they performed in Clutch over the last 90 days.


Using a 'Business Intelligence' tool, we can take this data and create nice charts to examine usage over time. We like to look at unique users per month and week and actions performed over the last 14 days and 90 days. From these charts, we can easily see which workflows are the most frequently used and how often Clutch is being used.




If you don’t have a charting tool, you’ll need a query for each of these charts.


This query will give you a count of unique users per month for the last 90 days.


This query will give you all actions performed over the last 14 days and how often they happened.


That’s it! Now, for some nuance. Our audit log only contains events sent to the middleware, so things like button clicks do not get logged. However, without installing tools like Google Analytics, mining the audit log is an easy way to see what your users are taking actions on.

If you are concerned about logging too much or logging sensitive information (e.g. PII), the Clutch middleware supports annotations that prevents a field from being logged.

Annotate the proto with “[ (clutch.api.v1.log) = false ]” to prevent a field from being logged.
	
At Lyft, for some actions we log entire response objects and are able to generate reports detailing information like how often our Slack integration is being utilized or which specific resources are being actioned against most often, all through the audit logs.

If you are interested in turning on Security Audit logging with Clutch, take a look at our guide.

Have a question? Join us in our Slack!


