describe("Landing Page", () => {
  before(() => {
    cy.visit("localhost:3000");
  });

  it("renders", () => {
    cy.get("#landing").should("be.visible");
  });
});
