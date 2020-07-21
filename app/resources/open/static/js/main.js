const {dialog} = require('electron').remote
const {app} = require('electron')

let main = {
  init: function () {
    // Wait for astilectron to be ready
    document.addEventListener('astilectron-ready', function () {
      // Listen
      main.listen();

      app.on('activate', (event, webContents, url, list, callback) => {
        astilectron.sendMessage("activateactivate", function (message) {

        });
      })
    });
  },
  open: function () {
    // Create message
    let message = {"name": "open"};

    dialog.showOpenDialog().then(result => {
      if (!result.canceled) {
        message.payload = result.filePaths[0];

        astilectron.sendMessage(message, function (message) {
          // Check error
          if (message.name === "error") {
            return
          }
        })
      }
    }).catch(err => {
      console.log(err)
    })
  },
  listen: function () {
    astilectron.onMessage(function (message) {
      switch (message.name) {
        case "menu-open":
          main.open();
      }
    });
  }
};
// export {main};
