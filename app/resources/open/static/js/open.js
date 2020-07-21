let open = {
  addFolder() {
    let openButton = document.createElement("button");
    openButton.onclick = function () {
      main.open()
    };
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
  init: function () {
    open.addFolder();

    // Wait for astilectron to be ready
    document.addEventListener('astilectron-ready', function () {
      // Listen
      open.listen();
    });
  },
  listen: function () {
    astilectron.onMessage(function (message) {
      switch (message.name) {
          case "menu-open":
              return {payload: "payload"};
      }
    });
  }
};
