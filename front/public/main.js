console.log('loaded')
document.addEventListener("wheel", function (e) {
    console.log("wheel",e.ctrlKey, e.deltaY )
    e.preventDefault();
}, { passive: false });


const urlFragment = window.location.hash;
let accessToken;

if (urlFragment) {
    if (urlFragment.includes("access_token")) {
        accessToken = urlFragment.split("access_token=")[1].split("&")[0];
        console.log("Access Token:", accessToken);
        document.querySelector('#login-button').style.display = 'none';

        disaplyContent();
    } else if (urlFragment.includes("error")) {
        const error = urlFragment.split("error=")[1].split("&")[0];
        console.error("OAuth 2.0 flow error:", error);
    }
}



function oauthSignIn() {
    // Google's OAuth 2.0 endpoint for requesting an access token
    var oauth2Endpoint = 'https://accounts.google.com/o/oauth2/v2/auth';

    // Create <form> element to submit parameters to OAuth 2.0 endpoint.
    var form = document.createElement('form');
    form.setAttribute('method', 'GET'); // Send as a GET request.
    form.setAttribute('action', oauth2Endpoint);

    // Parameters to pass to OAuth 2.0 endpoint.
    var params = {
        'client_id': '949757780668-mptrdgdc2hfvdu5bul8t7boog88nbd07.apps.googleusercontent.com',
        'redirect_uri': 'https://literate-space-invention-w9xw4xg5p6pc5x97-8080.app.github.dev/',
        'response_type': 'token',
        'scope': 'https://www.googleapis.com/auth/photoslibrary.readonly https://www.googleapis.com/auth/photospicker.mediaitems.readonly',
        'include_granted_scopes': 'true',
        'state': 'pass-through value'
    };

    // Add form parameters as hidden input values.
    for (var p in params) {
        var input = document.createElement('input');
        input.setAttribute('type', 'hidden');
        input.setAttribute('name', p);
        input.setAttribute('value', params[p]);
        form.appendChild(input);
    }

    // Add form to page and submit it to open the OAuth 2.0 endpoint.
    document.body.appendChild(form);
    form.submit();
}

function disaplyContent() {
    const base = 'https://photospicker.googleapis.com';

    fetch(base + '/v1/sessions', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + accessToken
        },
        json: true
    }).then((response) => response.json())
        .then((responseData) => {
            console.log('responseData Post', responseData);
            
  fetch(`${base}/v1/sessions/${responseData}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + req.user.token
    },
    json: true
  }).then((response) => response.json())
  .then((responseData) => {

    console.log('responseData get', responseData)
  })
        })


}
