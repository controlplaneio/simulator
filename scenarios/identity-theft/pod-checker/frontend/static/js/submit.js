async function submitForm() {
  try {
    const response = await fetch('/schema', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: document.getElementById('yaml').value,
    });
    const text = await response.text();
    console.log(text);
    var status = document.getElementById("icon")
    if (text.includes(Boolean(true))) {
      status.classList.add("fa-circle-check")
      status.classList.remove("fa-circle")
      status.classList.remove("fa-circle-xmark");
      document.getElementById('example').innerHTML = "";
    }
    else if (text.includes(Boolean(false))) {
      status.classList.add("fa-circle-xmark")
      status.classList.remove("fa-circle-check")
      status.classList.remove("fa-circle");
      getExample();
    }
    else {
      status.classList.add("fa-circle");
    }
  } catch (error) {
    console.error(error);
  }
}

async function getExample() {
  try {
    const response = await fetch('/example', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: '{"kind":"pod"}',
    });
    const text = await response.text();
    remove = text.replace(/["]+/g, "");
    yaml = remove.replace(/\\n/g, "&#10;");
    document.getElementById('example').innerHTML = yaml;
  } catch (error) {
    console.error(error);
  }

}