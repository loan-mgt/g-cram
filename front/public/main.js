console.log("loaded");

const VAPID_PUBLIC_KEY =
  "BO3EJqTUc89-9Q7lhTk3ypUyLZ2Lp7l8LTk6XVWyl9HD3DhP8JE5SUAsoHKZqGbw2766xV5oM3ixJ-0UhRkLD5E";
const base = "https://photospicker.googleapis.com";
const api = "http://localhost:8090/api/v1";
const mediaHolder = document.querySelector("#media-holder");
const loginButton = document.querySelector("#login-button");
const urlFragment = window.location.hash;
const videoCount = document.querySelector("#video-count");
let sessionId;
let pickerUri;
let picker;
let mediaItems;
let pullForImagesTimeout;

function handlePciker() {
  if (!picker && sessionId) {
    picker = window.open(pickerUri, "popup", "width=600,height=600");
    pullForImages();
  } else {
    picker.focus();
  }
}

// handle ongoing OAuth 2.0 flow
const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has("code")) {
  const code = urlParams.get("code");

  initUser(code);
  window.history.replaceState({}, "", "/");
} else if (urlParams.has("error")) {
  const error = urlParams.get("error");
  if (error) {
    console.error("OAuth 2.0 flow error:", error);
  }
} else {
  // currently logged in
  if (localStorage.getItem("name")) {
    getCurrentUser();
  }
  setScreen(0);

}

// user is loged in
function logedIn() {
  startWebSocket();
  disaplyContent();
  displayUserJobs();
}

function storeUserData(data) {
  console.log("User response:", data);
  if (data.userId !== undefined) localStorage.setItem("user_id", data.userId);
  if (data.id !== undefined) localStorage.setItem("user_id", data.id);
  if (data.userName !== undefined) localStorage.setItem("name", data.userName);
  if (data.name !== undefined) localStorage.setItem("name", data.name);
  if (data.picture !== undefined) localStorage.setItem("picture", data.picture);
  if (data.email !== undefined) localStorage.setItem("email", data.email);
  if (data.accessToken !== undefined)
    localStorage.setItem("access_token", data.accessToken);
}

function initUser(code) {
  fetch(api + "/user", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      code: code,
    }),
  })
    .then((response) => response.json())
    .then((data) => {
      storeUserData(data);
      logedIn();
    })
    .catch((error) => {
      console.error("Error initializing user:", error);
    });
}

async function getCurrentUser() {
  return fetch(api + "/user", {
    method: "GET",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Error initializing user");
      }
    })
    .then((data) => {
      storeUserData(data);
      logedIn();
    })
    .catch((error) => {
      console.error("Error initializing user:", error);
    });
}

function startWebSocket() {
  const ws = new WebSocket(
    `ws://localhost:8090/api/v1/ws?token=${localStorage.getItem("access_token")}&id=${localStorage.getItem("user_id")}`,
  );

  ws.onmessage = function (event) {
    const data = JSON.parse(event.data);
    console.log(data);
  };
}

function oauthSignIn() {
  // Google's OAuth 2.0 endpoint for requesting an access token
  var oauth2Endpoint = "https://accounts.google.com/o/oauth2/v2/auth";

  // Create <form> element to submit parameters to OAuth 2.0 endpoint.
  var form = document.createElement("form");
  form.setAttribute("method", "GET"); // Send as a GET request.
  form.setAttribute("action", oauth2Endpoint);

  // Parameters to pass to OAuth 2.0 endpoint.
  var params = {
    client_id:
      "949757780668-mptrdgdc2hfvdu5bul8t7boog88nbd07.apps.googleusercontent.com",
    redirect_uri: "http://localhost:8080/",
    response_type: "code",
    access_type: "offline",
    prompt: "consent", // remove long term
    scope:
      "https://www.googleapis.com/auth/photoslibrary.appendonly https://www.googleapis.com/auth/photospicker.mediaitems.readonly",
    include_granted_scopes: "true",
    state: "pass-through value",
  };

  let userMail = localStorage.getItem("email");
  if (userMail !== null && userMail !== undefined) {
    params.login_hint = userMail;
  }

  // Add form parameters as hidden input values.
  for (var p in params) {
    var input = document.createElement("input");
    input.setAttribute("type", "hidden");
    input.setAttribute("name", p);
    input.setAttribute("value", params[p]);
    form.appendChild(input);
  }

  // Add form to page and submit it to open the OAuth 2.0 endpoint.
  document.body.appendChild(form);
  form.submit();
}

function disaplyContent() {
  fetch(base + "/v1/sessions", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + localStorage.getItem("access_token"),
    },
    json: true,
  })
    .then((response) => response.json())
    .then((responseData) => {
      pickerUri = responseData.pickerUri;
      sessionId = responseData.id;
    });
}

function pullForImages() {
  fetch(`${base}/v1/sessions/${sessionId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + localStorage.getItem("access_token"),
    },
    json: true,
  })
    .then((response) => response.json())
    .then((responseData) => {
      if (responseData.mediaItemsSet) {
        fetchMediaItems(sessionId, localStorage.getItem("access_token"));
      } else {
        pullForImagesTimeout = setTimeout(() => pullForImages(), 5000);
      }
    });
}

function fetchMediaItems(id, token, size = 25) {
  let itemsQuery = `sessionId=${id}&pageSize=${size}`;

  const response = fetch(
    `https://photospicker.googleapis.com/v1/mediaItems?${itemsQuery}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
      json: true,
    },
  )
    .then((response) => response.json())
    .then((responseData) => {
      console.log("responseData", responseData);
      mediaItems = responseData.mediaItems;
      handleMedia(
        responseData.mediaItems.filter(
          (mediaItem) => mediaItem.type === "VIDEO",
        ),
      );
    });
}

function handleMedia(mediaItems) {
  const payload = mediaItems.map((m) => ({
    media_id: m.id,
    creation_date: Date.parse(m.createTime),
    filename: m.mediaFile.filename,
    base_url: m.mediaFile.baseUrl,
  }));

  fetch(`${api}/media`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(payload),
  })
    .then((response) => response.json())
    .then((responseData) => {
      setScreen(3);
      videoCount.textContent = responseData.nb_media;
    });
}

function imageFactory(mediaItem, w = 128, h = 128) {
  if (mediaItem.type === "VIDEO") {
    const video = document.createElement("video");
    video.id = mediaItem.id;
    video.width = w;
    video.height = h;
    video.referrerPolicy = "no-referrer";
    video.controls = true;
    loadVideoIntoVideo(video.id, mediaItem.mediaFile.baseUrl);
    return video;
  } else {
    const image = document.createElement("img");
    image.id = mediaItem.id;
    image.width = w;
    image.height = h;
    image.referrerPolicy = "no-referrer";
    loadImageIntoImg(
      image.id,
      mediaItem.mediaFile.baseUrl + "=w" + w + "-h" + h,
    );
    return image;
  }
}

const loadImageIntoImg = (imgId, baseUrl) => {
  fetch("http://localhost:8090/api/v1/get-image", {
    method: "POST",
    body: JSON.stringify({ baseUrl, id: imgId }),
    headers: {
      Authorization: "Bearer " + localStorage.getItem("access_token"),
    },
    json: true,
  }).then((res) => {
    res.blob().then((blob) => {
      document.getElementById(imgId).src = URL.createObjectURL(blob);
    });
  });
};

const loadVideoIntoVideo = (videoId, baseUrl) => {
  fetch("http://localhost:8090/api/v1/get-video", {
    method: "POST",
    body: JSON.stringify({ baseUrl, id: videoId }),
    headers: {
      Authorization: "Bearer " + localStorage.getItem("access_token"),
    },
    json: true,
  }).then((res) => {
    res.blob().then((blob) => {
      document.getElementById(videoId).src = URL.createObjectURL(blob);
    });
  });
};

function startCompression() {
  requestNotification();
  registerServiceWorker();
  fetch(`${api}/start`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(
      mediaItems
        .filter((mediaItem) => mediaItem.type === "VIDEO")
        .map((mediaItem) => ({
          id: mediaItem.id,
          creationDate: mediaItem.createTime,
          name: mediaItem.mediaFile.filename,
        })),
    ),
    json: true,
  })
    .then((response) => response.json())
    .then((_) => {
      setScreen(4);
    });
}

function requestNotification() {
  if (!("Notification" in window)) {
    // Check if the browser supports notifications
    alert("This browser does not support desktop notification");
  } else {
    Notification.requestPermission().then((permission) => {
      console.log("Notification", permission);
    });
  }
}

async function registerServiceWorker() {
  if (!localStorage.getItem("user_id")) {
    return;
  }

  // Check if service workers and push are supported
  if (!("serviceWorker" in navigator) || !("PushManager" in window)) {
    console.error("Browser does not support service workers or push messages.");
    return;
  }

  try {
    // Register the service worker
    const serviceWorkerRegistration = await navigator.serviceWorker.register(
      "./public/service-worker.js",
    );
    console.info("Service worker was registered.");
    console.info({ serviceWorkerRegistration });
    let registration = serviceWorkerRegistration;

    console.log(navigator.serviceWorker.controller, registration.pushManager);
    if (!registration.pushManager) {
      registration = await navigator.serviceWorker.ready;
    }

    // Make sure the registration is complete before proceeding

    // Check if already subscribed
    const subscribed = await registration.pushManager.getSubscription();
    const subscription = subscribed || (await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: VAPID_PUBLIC_KEY,
    }));

    // Send the subscription to the server
    await fetch(`${api}/user/subscription`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(subscription),
    });

    console.info("User subscription completed successfully.");
  } catch (error) {
    console.error(
      "Error during service worker registration or push subscription:",
      error,
    );
  }
}

async function unsubscribeButtonHandler() {
  // TODO
  const registration = await navigator.serviceWorker.getRegistration();
  const subscription = await registration.pushManager.getSubscription();
  fetch(`${api}/user/${localStorage.getItem("user_id")}/subscription`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ endpoint: subscription.endpoint }),
  });
  const unsubscribed = await subscription.unsubscribe();
  if (unsubscribed) {
    console.info("Successfully unsubscribed from push notifications.");
    unsubscribeButton.disabled = true;
    subscribeButton.disabled = false;
    notifyMeButton.disabled = true;
  }
}

function logout() {
  localStorage.clear();
  location.reload();
}


function setScreen(nb, load = false) {
  let screens = document.querySelectorAll('main');


  if (load){
    displayUserJobs();
  }


  screens.forEach((screen) => {
    if (screen.id === `screen-${nb}`) {
      screen.classList.remove('hidden')
      screen.classList.add('flex')
    } else {
      screen.classList.remove('flex')
      screen.classList.add('hidden')
    }
  })
}

function start() {
  if (localStorage.getItem("access_token")) {
    setScreen(2)
  } else {
    setScreen(1)
  }
}

function displayUserJobs() {
  fetch(api + "/job", {
    method: "GET",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => response.json())
    .then((jobs) => {
      const jobsList = document.querySelector("#jobs ul");
      jobsList.innerHTML = ""; // clear existing jobs

      jobs.forEach((job) => {
        const template = document.querySelector("#media-item-template");
        const jobItem = template.content.cloneNode(true);

        const timeSpan = jobItem.querySelector(".time");
        const time = formatTimeElapsed(job.timestamp);
        timeSpan.textContent = time;

        const sizeSpan = jobItem.querySelector(".size");
        const sizeSaved = formatFileSize(job.old_size - job.new_size);
        sizeSpan.textContent = sizeSaved;

         const img = jobItem.querySelector("img");
        if (job.nb_media === job.nb_media_done) {
          img.src = "/done.svg";
        } else {
          img.src = "/not-done.svg";
        }

        jobsList.appendChild(jobItem);
      });
    })
    .catch((error) => console.error("Error fetching jobs:", error));
}

function formatTimeElapsed(unixTime) {
  // Convert to milliseconds if it's in seconds (Unix timestamps less than 2^40)
  // This handles both millisecond and second formats
  if (unixTime < 2**40) {
    unixTime *= 1000;
  }
  
  const now = Date.now();
  const diffMs = now - unixTime;
  
  // Convert to various time units
  const seconds = Math.floor(diffMs / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const months = Math.floor(days / 30);
  const years = Math.floor(days / 365);
  
  // Return the most appropriate time format
  if (years > 0) {
    return `${years} ${years === 1 ? 'year' : 'years'}`;
  } else if (months > 0) {
    return `${months} ${months === 1 ? 'month' : 'months'}`;
  } else if (days > 0) {
    return `${days} ${days === 1 ? 'day' : 'days'}`;
  } else if (hours > 0) {
    return `${hours}h`;
  } else if (minutes > 0) {
    return `${minutes}min`;
  } else {
    return `${seconds}s`;
  }
}

function formatFileSize(bytes) {
  // Handle invalid input
  if (bytes === null || bytes === undefined || isNaN(bytes) || bytes < 0) {
    return "0 B";
  }
  
  // Array of units from bytes to petabytes
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  
  // Base for conversion (1024 for binary/IEC, 1000 for decimal/SI)
  const base = 1024;
  
  // If it's 0 bytes, return immediately
  if (bytes === 0) return `0 ${units[0]}`;
  
  // Calculate the appropriate unit index by taking the log of the file size
  const unitIndex = Math.floor(Math.log(bytes) / Math.log(base));
  
  // Make sure we don't exceed the available units
  const safeUnitIndex = Math.min(unitIndex, units.length - 1);
  
  // Calculate the size in the determined unit and round to integer
  const size = Math.floor(bytes / Math.pow(base, safeUnitIndex));
  
  // Return the formatted string
  return `${size} ${units[safeUnitIndex]}`;
}
