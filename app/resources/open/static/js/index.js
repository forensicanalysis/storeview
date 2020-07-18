const { dialog } = require('electron').remote

let index = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    addFolder() {
        let openButton = document.createElement("button");
        openButton.onclick = function() { index.explore() };
        openButton.innerHTML = `<i class="fa fa-folder"></i><span>Open forensicstore</span>`;
        document.getElementById("menu").appendChild(openButton);

      let newButton = document.createElement("button");
      newButton.className = "inactive";
      // newButton.onclick = function() { index.explore() };
      newButton.innerHTML = `<i class="fa fa-plus"></i><span>New forensicstore</span>`;
      document.getElementById("menu").appendChild(newButton);

      let imageButton = document.createElement("button");
      imageButton.className = "inactive";
      // button.onclick = function() { index.explore() };
      imageButton.innerHTML = `<i class="fa fa-hdd-o"></i><span>Import disk image</span>`;
      document.getElementById("menu").appendChild(imageButton);
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
        });
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
            }
        });
    }
};
