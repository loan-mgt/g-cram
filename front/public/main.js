console.log('loaded')

const VAPID_PUBLIC_KEY = 'BO3EJqTUc89-9Q7lhTk3ypUyLZ2Lp7l8LTk6XVWyl9HD3DhP8JE5SUAsoHKZqGbw2766xV5oM3ixJ-0UhRkLD5E';
const base = 'https://photospicker.googleapis.com';
const api = 'http://localhost:8090/api/v1';
const mediaHolder = document.querySelector('#media-holder');
const loginButton = document.querySelector('#login-button');
const pickerLink = document.querySelector('#picker');
const pickerButton = document.querySelector('#picker2');
const urlFragment = window.location.hash;
let sessionId;
let pickerUri;
let picker;
let mediaItems;
pickerButton.disable = true;


pickerButton.addEventListener('click', function () {
    if (!picker) {
        picker = window.open(pickerUri, 'popup', 'width=600,height=600');
        pullForImages();
    } else {
        picker.focus();
    }
});


const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('code')) {
    const code = urlParams.get('code');

    initUser(code);
    window.history.replaceState({}, '', '/');
    logedIn();
} else if (urlParams.has('error')) {
    const error = urlParams.get('error');
    if (error) {
        console.error("OAuth 2.0 flow error:", error);
    }
}

if (localStorage.getItem('userId')) {
    getAccessToken()
        .then(_ => {
            logedIn();
        })
}

function logedIn() {
    loginButton.style.display = 'none';
    startWebSocket();
    disaplyContent();
    getUserInfo();
}

function initUser(code) {
    fetch(api + '/user', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            code: code,
        })
    })
        .then(response => response.json())
        .then(data => {
            console.log('User response:', data);
            localStorage.setItem('user_id', data.userId);
            localStorage.setItem('user_name', data.userName);
            localStorage.setItem('user_mail', data.userMail);
            localStorage.setItem('access_token', data.accessToken);
        })
        .catch(error => {
            console.error('Error initializing user:', error);
        });
}

function startWebSocket() {
    const clientId = 'bob'

    const ws = new WebSocket(`ws://localhost:8090/api/v1/ws?token=${localStorage.getItem('access_token')}&id=${localStorage.getItem('user_id')}`);

    ws.onmessage = function (event) {
        const data = JSON.parse(event.data);
        console.log(data);
    };
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
        'response_type': 'code',
        'access_type': 'offline',
        'scope': 'https://www.googleapis.com/auth/photoslibrary https://www.googleapis.com/auth/photospicker.mediaitems.readonly',
        'include_granted_scopes': 'true',
        'state': 'pass-through value',
        'login_hint': localStorage.getItem('user_mail')
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

    fetch(base + '/v1/sessions', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('access_token')
        },
        json: true,
    }).then((response) => response.json())
        .then((responseData) => {
            console.log('responseData Post', responseData);
            pickerUri = responseData.pickerUri
            pickerButton.disable = false;

            pickerLink.href = responseData.pickerUri;
            sessionId = responseData.id;



        })
}

function pullForImages() {
    fetch(`${base}/v1/sessions/${sessionId}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('access_token')
        },
        json: true
    }).then((response) => response.json())
        .then((responseData) => {

            console.log('responseData get', responseData)
            if (responseData.mediaItemsSet) {
                console.log("ended", responseData.mediaItemsSet);
                fetchMediaItems(sessionId, localStorage.getItem('access_token'));
            } else {
                setTimeout(() => pullForImages(), 5000);
            }
        })

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
            mediaItems = responseData.mediaItems;
            readerMedia(responseData.mediaItems.filter(mediaItem => mediaItem.type === 'VIDEO'));

        });
}

function readerMedia(mediaItems) {
    mediaItems.forEach(mediaItem => {
        mediaHolder.appendChild(imageFactory(mediaItem));
    })
}

function imageFactory(mediaItem, w = 128, h = 128) {

    if (mediaItem.type === 'VIDEO') {
        const video = document.createElement('video');
        video.id = mediaItem.id
        video.width = w;
        video.height = h;
        video.referrerPolicy = "no-referrer";
        video.controls = true;
        loadVideoIntoVideo(video.id, mediaItem.mediaFile.baseUrl);
        return video;

    } else {
        const image = document.createElement('img');
        image.id = mediaItem.id
        image.width = w;
        image.height = h;
        image.referrerPolicy = "no-referrer";
        loadImageIntoImg(image.id, mediaItem.mediaFile.baseUrl + "=w" + w + "-h" + h);
        return image;

    }
}


const loadImageIntoImg = (imgId, baseUrl) => {
    fetch('http://localhost:8090/api/v1/get-image', {
        method: 'POST',
        body: JSON.stringify({ baseUrl, id: imgId }),
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        },
        json: true,
    }).then(res => {
        res.blob().then(blob => {
            document.getElementById(imgId).src = URL.createObjectURL(blob)
        })
    })
}

const loadVideoIntoVideo = (videoId, baseUrl) => {
    fetch('http://localhost:8090/api/v1/get-video', {
        method: 'POST',
        body: JSON.stringify({ baseUrl, id: videoId }),
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        },
        json: true,
    }).then(res => {
        res.blob().then(blob => {
            document.getElementById(videoId).src = URL.createObjectURL(blob)
        })
    })
}


function startCompression() {
    requestNotification();
    registerServiceWorker();
    fetch(`${api}/start?id=${localStorage.getItem('user_id')}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('access_token')
        },
        body: JSON.stringify(mediaItems.filter(mediaItem => mediaItem.type === 'VIDEO')
            .map(mediaItem => ({
                id: mediaItem.id,
                creationDate: mediaItem.createTime,
                name: mediaItem.mediaFile.filename
            }))),
        json: true
    }).then((response) => response.json())
        .then((responseData) => {
            console.log('responseData', responseData)
        });
}
function loginSuccess(response) {
    console.log('loginSuccess')
    const responsePayload = decodeJwtResponse(response.credential);

    console.log("ID: " + responsePayload.sub);
    console.log('Full Name: ' + responsePayload.name);
    console.log('Given Name: ' + responsePayload.given_name);
    console.log('Family Name: ' + responsePayload.family_name);
    console.log("Image URL: " + responsePayload.picture);
    console.log("Email: " + responsePayload.email);

}
function decodeJwtResponse(token) {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

function getUserInfo() {
    fetch('https://www.googleapis.com/oauth2/v2/userinfo', {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('access_token')}`
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('User info:', data);
        })
        .catch(error => {
            console.error('Error fetching user info:', error);
        });
}

function requestNotification() {
    if (!("Notification" in window)) {
        // Check if the browser supports notifications
        alert("This browser does not support desktop notification");
    } else {

        Notification.requestPermission()
            .then((permission) => {
                console.log("Notification", permission);
                if (permission === "granted") {
                    new Notification("G-CRAM", {
                        body: "Starting Compression",
                    });
                }
            })
    }
}

async function registerServiceWorker() {
    if (!localStorage.getItem('user_id')) {
        return
    }

    // Check if service workers and push are supported
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
        console.error('Browser does not support service workers or push messages.');
        return;
    }

    try {
        // Register the service worker
        const serviceWorkerRegistration = await navigator.serviceWorker.register('./public/service-worker.js');
        console.info('Service worker was registered.');
        console.info({ serviceWorkerRegistration });
        let registration = serviceWorkerRegistration;

        console.log(navigator.serviceWorker.controller, registration.pushManager)
        if (!registration.pushManager) {
            registration = await navigator.serviceWorker.ready;
        }

        // Make sure the registration is complete before proceeding

        // Check if already subscribed
        const subscribed = await registration.pushManager.getSubscription();
        if (subscribed) {
            console.info('User is already subscribed.');
            registration.update();
            return;
        }

        // Subscribe the user
        const subscription = await registration.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: VAPID_PUBLIC_KEY
        });

        // Send the subscription to the server
        await fetch(`${api}/user/${localStorage.getItem('user_id')}/subscription`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(subscription)
        });

        console.info('User subscription completed successfully.');
    } catch (error) {
        console.error('Error during service worker registration or push subscription:', error);
    }
}

async function unsubscribeButtonHandler() {
    // TODO
    const registration = await navigator.serviceWorker.getRegistration();
    const subscription = await registration.pushManager.getSubscription();
    fetch(`${api}/user/${localStorage.getItem('user_id')}/subscription`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ endpoint: subscription.endpoint })
    });
    const unsubscribed = await subscription.unsubscribe();
    if (unsubscribed) {
        console.info('Successfully unsubscribed from push notifications.');
        unsubscribeButton.disabled = true;
        subscribeButton.disabled = false;
        notifyMeButton.disabled = true;
    }
}
