import React from "react";
import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom/extend-expect";
import "@testing-library/jest-dom";
import MarkdownPane from "../markdown-pane";
const mdContent = `
## A subtitle

- Meow meow meow
- a bulleted line of text
`

test("Renders a title and content", () => {
    const wrapper = render(<MarkdownPane onClose={() => {}} title={"The Title"} markdownText={mdContent} />)

    // expect the dashes to be converted, so they won't be found
    expect(screen.getByText("-")).toBeNull();
    expect(screen.getByText("The Title")).toBeInTheDocument();
    expect(screen.getByText("Meow meow meow")).toBeInTheDocument();
})

