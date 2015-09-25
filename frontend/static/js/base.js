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

/** Member API namespace. */
osa.api.member = osa.api.member || {};

/** UI namespace. */
osa.api.ui = osa.api.ui || {};

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
      document.querySelector('#updateMember').disabled = false;
      osa.api.member.current();
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
 * Clears auth token and resets UI when user logs out.
 */
osa.api.signout = function() {
  gapi.auth.setToken(null);
  gapi.auth.signOut();
  osa.api.signedIn = false;
  document.querySelector('#signinButton').textContent = 'Sign in';
  document.querySelector('#addEvent').disabled = true;
  document.querySelector('#updateMember').disabled = true;
  osa.api.member.print();
  osa.api.ui.cleanLog();
}

/**
 * Presents the user with the authorization popup when signing in.
 * Cleans up when user signs out.
 */
osa.api.auth = function() {
  if (!osa.api.signedIn) { // Sign in
    osa.api.signin(false, osa.api.userAuthed);
  } else { // Sign out
    osa.api.signout();
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
osa.api.ui.msg = function(msg) {
  var element = document.createElement('div');
  element.classList.add('p');
  element.innerHTML = msg;
  document.querySelector('#outputLog').appendChild(element);
};

/**
 * Deletes all children of a specified node
 * @param {Object} node The node to delete children of
 */
osa.api.ui.clean = function(node) {
  var fc = node.firstChild;
  while( fc ) {
      node.removeChild( fc );
      fc = node.firstChild;
  }
}

/**
 * Removes all items in the output log
 */
osa.api.ui.cleanLog = function() {
  osa.api.ui.clean(document.querySelector('#outputLog'))
}

/**
 * Lists events via the API.
 */
osa.api.events.list = function() {
  gapi.client.events.events.list().execute(
    function(resp) {
      osa.api.ui.cleanLog();
      if (!resp.code) {
        resp.events = resp.events || [];
        if (resp.events.length == 0) {
          osa.api.ui.msg('No events found.');
          return;
        }
        for (var i = 0; i < resp.events.length; i++) {
          osa.api.events.print(resp.events[i]);
        }
      } else {
        osa.api.ui.msg('Error: ' + resp.message);
      }
    }
  );
};

/**
 * Adds a new event.
 * @param {Object} event the even object to add
 */
osa.api.events.add = function(event) {
  gapi.client.events.events.add(event).execute(
    function(resp) {
      osa.api.ui.cleanLog();
      if (!resp.code) {
        osa.api.ui.msg('Successively added event!');
      } else {
        osa.api.ui.msg('Error: ' + resp.message);
      }
    });
};

/**
 * Fetches the user's profile.
 */
osa.api.member.current = function() {
  gapi.client.member.member.current().execute(
    function(resp) {
      osa.api.member.print(resp);
    }
  );
};

/**
 * Updates the user's profile with the member object
 * @param {Object} member the member object to update with
 */
osa.api.member.update = function(member) {
  gapi.client.member.member.update(member).execute(
    function(resp) {
      osa.api.ui.cleanLog();
      if (!resp.code) {
        osa.api.ui.msg('Successively updated profile!');
      } else {
        osa.api.ui.msg('Error: ' + resp.message);
      }
    }
  );
};

/**
 * Loads the resp as a member object if it is valid
 * @param {Object} resp the member object or error response
 */
osa.api.member.print = function(resp) {
  var first = document.querySelector('#firstName');
  var last = document.querySelector('#lastName');
  var rel = document.querySelector('#relationship');
  if (resp && !resp.code) {
    first.value = resp.name.first;
    last.value = resp.name.last;
    rel.value = resp.relationship;
  } else {
    first.value = '';
    last.value = '';
    rel.value = '';
  }
}

/**
 * Enables the button callbacks in the UI.
 */
osa.api.ui.enableButtons = function() {
  var listEvents = document.querySelector('#listEvents');
  listEvents.addEventListener('click',
      osa.api.events.list);

  var addEvent = document.querySelector('#addEvent');
  addEvent.addEventListener('click', function() {
    osa.api.events.add({
      'title': document.querySelector('#eventTitle').value,
      'description': document.querySelector('#description').value
    });
  });

  var updateMember = document.querySelector('#updateMember');
  updateMember.addEventListener('click', function() {
    osa.api.member.update({
      'relationship': document.querySelector('#relationship').value,
      'name': {
        'first': document.querySelector('#firstName').value,
        'last': document.querySelector('#lastName').value
      }
    });
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
      osa.api.ui.enableButtons();
      osa.api.signin(true, osa.api.userAuthed);
    }
  }

  apisToLoad = 3; // must match number of calls to gapi.client.load()
  gapi.client.load('events', 'v1', callback, apiRoot);
  gapi.client.load('member', 'v1', callback, apiRoot);
  gapi.client.load('oauth2', 'v2', callback);
};
