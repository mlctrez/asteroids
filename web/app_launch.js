(() => {
  const go = new Go();
  fetch("app.wasm").then((response) => {
    if (!response.ok) {
      let errorPre = document.createElement("pre");
      document.body.appendChild(errorPre);
      response.text().then((text) => errorPre.innerText = text)
    } else {
      WebAssembly.instantiateStreaming(response, go.importObject).then(
        (result) => go.run(result.instance)
          .then(() => console.log("go.run exited"))
          .catch(err => console.log("error ", err))
      ).catch(err => console.log("error ", err))
    }
  })
})();