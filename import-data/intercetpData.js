// Save a reference to the original XMLHttpRequest object
const originalXHR = window.XMLHttpRequest;

// Array to store response texts
const responseTexts = [];

// Override the original XMLHttpRequest object with a custom implementation
window.XMLHttpRequest = function() {
    // Create a new instance of the original XMLHttpRequest object
    const xhr = new originalXHR();

    // Save a reference to the original open method
    const originalOpen = xhr.open;

    // Override the open method to intercept requests
    xhr.open = function(method, url) {
        // Check if the URL contains the specified string
        if (url.includes('https://solomonk.fr/ajax/get_spell_infos.php')) {
            // Save method and URL for logging
            const interceptedRequest = {
                method: method,
                url: url,
                requestBody: null
            };

            // Override the send method to intercept the request body
            const originalSend = xhr.send;
            xhr.send = function(body) {
                // Save the request body
                interceptedRequest.requestBody = body;

                // Call the original send method
                originalSend.apply(this, arguments);
            };

            // Log the intercepted request after it's sent
            xhr.addEventListener('load', function() {
                console.log('URL:', interceptedRequest.url);

                // Store response text in the array
                responseTexts.push(xhr.responseText);
            });
        }

        // Call the original open method
        originalOpen.apply(this, arguments);
    };

    // Return the modified XMLHttpRequest object
    return xhr;
};

// Now, XMLHttpRequests made to the specified URL will be intercepted, and their response texts will be stored in the array
console.log(responseTexts.join(',\n'))