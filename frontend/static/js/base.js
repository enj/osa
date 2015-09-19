/**
 * @fileoverview
 * Provides methods for the Events UI and interaction with the Events API.
 *
 * @author theenjeru@gmail.com (Monis Khan)
 *
 * Heavily based on example code.
 */

/** OSA global namespace. */
var osa = osa || {};

/** API namespace for OSA. */
osa.api = osa.api || {};

/** Events API namespace. */
osa.api.events = osa.api.events || {};

/**
 * Client ID of OSA app
 * @type {string}
 */
osa.api.CLIENT_ID =
    '206868860697-h39gavnuht6g1mle7esc0hva3euq33k6.apps.googleusercontent.com';

/**
 * Scopes used by OSA app
 * @type {string}
 */
osa.api.SCOPES =
    'https://www.googleapis.com/auth/userinfo.email';

/**
 * Whether or not the user is signed in.
 * @type {boolean}
 */
osa.api.signedIn = false;

/**
 * Loads the application UI after the user has completed auth.
 */
osa.api.userAuthed = function() {
  var request = gapi.client.oauth2.userinfo.get().execute(function(resp) {
    if (!resp.code) {
      osa.api.signedIn = true;
      document.querySelector('#signinButton').textContent = 'Sign out';
      document.querySelector('#addEvent').disabled = false;
    }
  });
};

/**
 * Handles the auth flow, with the given value for immediate mode.
 * @param {boolean} mode Whether or not to use immediate mode.
 * @param {Function} callback Callback to call on completion.
 */
osa.api.signin = function(mode, callback) {
  gapi.auth.authorize({client_id: osa.api.CLIENT_ID,
      scope: osa.api.SCOPES, immediate: mode},
      callback);
};

/**
 * Presents the user with the authorization popup.
 */
osa.api.auth = function() {
  if (!osa.api.signedIn) { // Sign in
    osa.api.signin(false, osa.api.userAuthed);
  } else { // Sign out
    osa.api.signedIn = false;
    document.querySelector('#signinButton').textContent = 'Sign in';
    document.querySelector('#addEvent').disabled = true;
    gapi.auth.setToken(null);
    gapi.auth.signOut();
  }
};

/**
 * Prints an event to the output log.
 * param {Object} event Event to print.
 */
osa.api.events.print = function(event) {
  var element = document.createElement('ul');
  element.classList.add('list-group');
  var child = document.createElement('li');
  child.classList.add('list-group-item');
  child.innerHTML = event.title + ': ' + event.description;
  element.appendChild(child);
  document.querySelector('#outputLog').appendChild(element);
};

/**
 * Prints a message to the output log.
 * param {string} msg Message to print
 */
osa.api.events.msg = function(msg) {
  var element = document.createElement('div');
  element.classList.add('p');
  element.innerHTML = msg;
  document.querySelector('#outputLog').appendChild(element);
};

/**
 * Deletes all children of a specified node
 * @param {Object} node The node to delete children of
 */
clean = function(node) {
  var fc = node.firstChild;
  while( fc ) {
      node.removeChild( fc );
      fc = node.firstChild;
  }
}

/**
 * Removes all items in the output log
 */
cleanLog = function() {
  clean(document.querySelector('#outputLog'))
}

/**
 * Lists events via the API.
 */
osa.api.events.listEvents = function() {
  gapi.client.events.events.list().execute(
      function(resp) {
        if (!resp.code) {
          cleanLog();
          resp.events = resp.events || [];
          for (var i = 0; i < resp.events.length; i++) {
            osa.api.events.print(resp.events[i]);
          }
        }
      });
};

/**
 * Adds a new event.
 * @param {string} title the title of the event
 * @param {string} count a description of the event
 */
osa.api.events.addEvent = function(title, description) {
  gapi.client.events.events.add({
      'title': title,
      'description': description
    }).execute(function(resp) {
      if (!resp.code) {
        cleanLog();
        osa.api.events.msg('Successively added event!');
      }
    });
};

/**
 * Enables the button callbacks in the UI.
 */
osa.api.events.enableButtons = function() {
  var listEvents = document.querySelector('#listEvents');
  listEvents.addEventListener('click',
      osa.api.events.listEvents);

  var addEvent = document.querySelector('#addEvent');
  addEvent.addEventListener('click', function() {
    osa.api.events.addEvent(
        document.querySelector('#eventTitle').value,
        document.querySelector('#description').value);
  });

  var signinButton = document.querySelector('#signinButton');
  signinButton.addEventListener('click', osa.api.auth);
};

/**
 * Initializes the application.
 * @param {string} apiRoot Root of the API's path.
 */
osa.api.init = function(apiRoot) {
  // Loads the OAuth and Events APIs asynchronously,
  // and triggers login when they have completed.
  var apisToLoad;
  var callback = function() {
    if (--apisToLoad == 0) {
      osa.api.events.enableButtons();
      osa.api.signin(true, osa.api.userAuthed);
    }
  }

  apisToLoad = 2; // must match number of calls to gapi.client.load()
  gapi.client.load('events', 'v1', callback, apiRoot);
  gapi.client.load('oauth2', 'v2', callback);
};
