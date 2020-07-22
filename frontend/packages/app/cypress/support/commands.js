Cypress.Commands.add("element", e2eId => {
  const ids = e2eId.split(" ");
  const locators = ids.map(id => {
    return `*[data-qa="${id}"]`;
  });
  return cy.get(locators.join(" "));
});

Cypress.Commands.add("descendent", { prevSubject: "true" }, (subject, e2eId) => {
  const ids = e2eId.split(" ");
  const locators = ids.map(id => {
    return `*[data-qa="${id}"]`;
  });
  return cy.wrap(subject).find(locators.join(" "));
});
