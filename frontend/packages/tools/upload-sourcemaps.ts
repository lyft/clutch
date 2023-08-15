const { browser } = require("@bugsnag/source-maps");

// interface browser  {
//   async function uploadOne (UploadSingleOpts): Promise<void>
//   async function uploadMultiple (UploadMultipleOpts): Promise<void>
// }

console.log("PROCESS", process.env.REACT_APP_BUGSNAG_API_TOKEN);

// browser.uploadOne({

// })
