import * as React from "react";
import styled from "@emotion/styled";

const HorizontalRuleBase = ({ children, ...props }: HorizontalRuleProps) => (
    <div {...props}>
        <div className="line">
            <span />
        </div>
        {React.Children.count(children) > 0 && <div className="content">{children}</div>}
        <div className="line">
            <span />
        </div>
    </div>
);

export type HorizontalRuleProps = {
    children: React.ReactNode
}

const StyledHorizontalRule = styled(HorizontalRuleBase)({
    alignItems: "center",
    display: "flex",
    flexDirection: "row",

    ".line": {
        flex: "1 1 auto"
    },

    ".line > span": {
        display: "block",
        borderTop: "1px solid rgba(13, 16, 48, 0.12)"
    },

    ".content": {
        padding: "0 30px",
        fontWeight: "bold",
        color: "rgba(13, 16, 48, 0.38)",
        textTransform: "uppercase",
        display: "inline-flex",
        alignItems: "center",
    },
});

export const HorizontalRule = ({ children }: HorizontalRuleProps) => <StyledHorizontalRule>{children}</StyledHorizontalRule>

export default HorizontalRule;