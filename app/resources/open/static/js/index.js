const { dialog } = require('electron').remote

let index = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    addFolder() {
        let button = document.createElement("button");
        button.onclick = function() { index.explore() };
        button.innerHTML = `<i class="fa fa-folder"></i><span>Open forensicstore</span>`;
        document.getElementById("dirs").appendChild(button)
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        index.addFolder();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
        })
    },
    explore: function() {
        // Create message
        let message = {"name": "open"};

        dialog.showOpenDialog().then(result => {
          if (!result.canceled) {
            // Send message

            message.payload = result.filePaths[0];

            asticode.loader.show();
            astilectron.sendMessage(message, function(message) {
              // Init
              asticode.loader.hide();

              // Check error
              if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
              }
            })
          }
        }).catch(err => {
          console.log(err)
        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    }
};
