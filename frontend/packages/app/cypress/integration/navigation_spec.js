const MENU_BUTTON = "menuBtn";
const DRAWER = "drawer";
const WORKFLOW_GROUP = "workflowGroup";
const TOGGLE = "toggle";

describe("Navigation drawer", () => {
  before(() => {
    cy.visit("localhost:3000");
    cy.element(MENU_BUTTON).click();
    cy.element(DRAWER).should("be.visible");
  });

  it("displays groups", () => {
    cy.element(WORKFLOW_GROUP).each(group => {
      cy.wrap(group).should("be.visible");
    });
  });

  it("displays and hides routes", () => {
    cy.element(WORKFLOW_GROUP).each((_, idx) => {
      cy.element(WORKFLOW_GROUP).eq(idx).descendent(TOGGLE).children().first().click();
      cy.element(WORKFLOW_GROUP)
        .eq(idx)
        .find("a")
        .each(link => {
          cy.wrap(link).should("have.attr", "href");
        });
      cy.element(WORKFLOW_GROUP).eq(idx).descendent(TOGGLE).children().first().click();
      cy.element(WORKFLOW_GROUP).eq(idx).find("a").should("not.be.visible");
    });
  });

  describe("routes to homepage", () => {
    it("via nav icon", () => {
      cy.element(DRAWER).element("logo").click();
      cy.url().should("equal", "http://localhost:3000/");
    });

    it("via nav title", () => {
      cy.element(MENU_BUTTON).click();
      cy.element(DRAWER).element("title").click();
      cy.url().should("equal", "http://localhost:3000/");
    });
  });

  describe("routes to workflows", () => {
    const groupItemId = "workflowGroupItem";
    beforeEach(() => {
      cy.element(MENU_BUTTON).click();
    });

    it("can route correctly", () => {
      return cy.element(WORKFLOW_GROUP).each((_, groupIdx) => {
        cy.element(WORKFLOW_GROUP).eq(groupIdx).descendent(TOGGLE).children().first().click();
        cy.element(WORKFLOW_GROUP)
          .eq(groupIdx)
          .descendent(groupItemId)
          .each((__, itemIdx) => {
            cy.element(WORKFLOW_GROUP)
              .eq(groupIdx)
              .descendent(groupItemId)
              .eq(itemIdx)
              .should("be.visible");
            cy.element(WORKFLOW_GROUP)
              .eq(groupIdx)
              .descendent(groupItemId)
              .eq(itemIdx)
              .should("have.attr", "href")
              .then(href => {
                cy.element(WORKFLOW_GROUP).eq(groupIdx).descendent(groupItemId).eq(itemIdx).click();
                cy.url().should("include", href);
                cy.element(MENU_BUTTON).click();
              });

            // TODO: validate header of workflow here when it's landed
          });
      });
    });
  });
});
