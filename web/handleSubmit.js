function handleSubmit(event, url) {
    event.preventDefault();
    const formElements = event.target.elements;
    const formData = {};
    for (let element of formElements) {
        if (element.tagName === 'INPUT' || element.tagName === 'SELECT') {
            formData[element.name || element.id || element.placeholder || element.previousSibling.textContent.trim()] = element.value;
        }
    }
    formData.operation = event.target.id
    console.log("Sucsess: ", formData);

    // Send the form data to the server
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        if (data.message !== "ok" && data.message !== "") {
            cleanedStr = data.message.replace(/\x00/g, '');
            alert(cleanedStr);
        } else {
            location.reload();
        }
    })
    .catch((error) => {
        console.error('Error:', error);
        alert('Error submitting the form: ' + error);
    });
}
