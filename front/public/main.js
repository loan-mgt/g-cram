console.log('loaded')

const mediaHolder = document.querySelector('#media-holder');
const loginButton = document.querySelector('#login-button');
const pickerLink = document.querySelector('#picker');
const pickerButton = document.querySelector('#picker2');
const urlFragment = window.location.hash;
let accessToken;
let sessionId;
let pickerUri;
let picker;
pickerButton.disable = true;

pickerButton.addEventListener('click', function() {
    if (!picker) {
        picker = window.open(pickerUri, 'popup', 'width=600,height=600');
    } else {
        picker.focus();
    }
});

if (urlFragment) {
    if (urlFragment.includes("access_token")) {
        accessToken = urlFragment.split("access_token=")[1].split("&")[0];
        console.log("Access Token:", accessToken);
        loginButton.style.display = 'none';

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
        'redirect_uri': 'http://localhost:8080/',
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
            pickerUri = responseData.pickerUri
            pickerButton.disable = false;

            pickerLink.href = responseData.pickerUri;
            sessionId = responseData.id;

            pullForImages();

        })

    function pullForImages() {
        fetch(`${base}/v1/sessions/${sessionId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + accessToken
            },
            json: true
        }).then((response) => response.json())
            .then((responseData) => {

                console.log('responseData get', responseData)
                if (responseData.mediaItemsSet) {
                    console.log("ended", responseData.mediaItemsSet);
                    fetchMediaItems(sessionId, accessToken);
                } else {
                    setTimeout(() => pullForImages(), 5000);
                }
            })

    }


}

function fetchMediaItems(id, token, size = 25) {
  
    let itemsQuery = `sessionId=${id}&pageSize=${size}`

    const response = fetch(`https://photospicker.googleapis.com/v1/mediaItems?${itemsQuery}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token
      },
      json: true
    }).then((response) => response.json())
    .then((responseData) => {
      console.log('responseData', responseData)

      readerMedia(responseData.mediaItems);

    });

    /**
     * [
    {
        "id": "ADXtojJYFMtZiW4004m5FSaXgmvTRr2RnP3gUpewIEyjsKKOVuovTlck7-xn4QmgLRiA3UHKhOIrBZERja6LmW3yVEutmw7TSQ",
        "createTime": "2025-02-27T23:34:05.466Z",
        "type": "PHOTO",
        "mediaFile": {
            "baseUrl": "https://lh3.googleusercontent.com/ppa/AGmLvJg2jwoR5-QWUP-pGpLyngWlYPNLmHKlhWGUcpi4q9JM2No__Oy4tdxXX0Tfh9C_mtBcW9w_SQ",
            "mimeType": "image/jpeg",
            "mediaFileMetadata": {
                "width": 3264,
                "height": 2448,
                "cameraMake": "Google",
                "cameraModel": "Pixel 4a (5G)",
                "photoMetadata": {
                    "focalLength": 2.57,
                    "apertureFNumber": 2,
                    "isoEquivalent": 1822,
                    "exposureTime": "0.055546998s"
                }
            },
            "filename": "PXL_20250227_233405466.jpg"
        }
    }
]
     */
}

function readerMedia(mediaItems){
    mediaItems.forEach(mediaItem => {
        mediaHolder.appendChild(imageFactory(mediaItem));
    })
}

function imageFactory(mediaItem, w = 128, h = 128) {
    const image = document.createElement('img');
    image.id = mediaItem.id
    image.width = w;
    image.height = h;
    image.referrerPolicy = "no-referrer";
    loadImageIntoImg(image.id, mediaItem.mediaFile.baseUrl+"=w"+w+"-h"+h);
    return image;
}


const loadImageIntoImg = (imgId, baseUrl) => {
    const response = fetch(baseUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + accessToken,
        'Referrer-Policy': 'no-referrer',
    },
      json: true,
     }).then(res => {
       res.blob().then(blob => {
         document.getElementById(imgId).src = URL.createObjectURL(blob)
       })
     })
  }