let brokerBtn = document.getElementById("brokerBtn");
let output = document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");

brokerBtn.addEventListener("click", () => {
  const body = {
    method: "POST",
  };

  fetch("http://localhost:8080", body)
    .then((response) => response.json())
    .then((data) => {
      sent.innerHTML("empty post request");
      received.innerHTML = JSON.stringify(data, undefined, 4);

      if (data.error) {
        console.log(data.error);
      } else {
        output.innerHTML += `<br><strong>Response from Broker service</strong>: ${data.message}`;
      }
    })
    .catch((err) => {
      output.innerHTML += `<br><br>Error: ${err}`;
    });
});
