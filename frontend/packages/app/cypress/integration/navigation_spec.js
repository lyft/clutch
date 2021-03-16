const DRAWER = "drawer";
const WORKFLOW_GROUP = "workflowGroup";
const WORKFLOW_GROUP_ITEM = "workflowGroupItem";

describe("Navigation drawer", () => {
  before(() => {
    cy.visit("localhost:3000");
    cy.element(DRAWER).should("be.visible");
  });

  it("displays groups", () => {
    cy.element(WORKFLOW_GROUP).each(group => {
      cy.wrap(group).should("be.visible");
    });
  });

  it("displays and hides routes", () => {
    cy.element(WORKFLOW_GROUP).each((_, idx) => {
      cy.element(WORKFLOW_GROUP).eq(idx).click();
      cy.element(WORKFLOW_GROUP_ITEM).each(link => {
        cy.wrap(link).should("have.attr", "href");
      });
      cy.element(WORKFLOW_GROUP).eq(idx).click();
      cy.element(WORKFLOW_GROUP).eq(idx).descendent(WORKFLOW_GROUP_ITEM).should("not.exist");
    });
  });

  describe("routes to workflows", () => {
    it("can route correctly", () => {
      return cy.element(WORKFLOW_GROUP).each((_, groupIdx) => {
        cy.element(WORKFLOW_GROUP).eq(groupIdx).click();
        cy.element(WORKFLOW_GROUP_ITEM).each((__, linkIdx) => {
          cy.element(WORKFLOW_GROUP_ITEM).eq(linkIdx).should("be.visible");
          cy.element(WORKFLOW_GROUP_ITEM)
            .eq(linkIdx)
            .should("have.attr", "href")
            .then(href => {
              cy.element(WORKFLOW_GROUP_ITEM).eq(linkIdx).click();
              cy.url().should("include", href);
            });
          cy.element(WORKFLOW_GROUP).eq(groupIdx).click();
          // TODO: validate header of workflow here when it's landed
        });
      });
    });
  });
});
