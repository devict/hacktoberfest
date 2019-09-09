// On document load we wire up all of the events
$(function() {
  if ($('#sponsorShareModal').length) {
    $('#allowSharing').click(allowSharing);
    $('#dismissModal').click(function() {setShareInfoState(false)});
    $('#sponsorShareModal').modal();
  } else {
    getShareInfo();
  }

  $('#issues_show').click(loadIssues);
  $('#check').click(checkPRs);
  $('#share_info').change(saveShareInfo);
});

function loadIssues() {
  var btn = $('#issues_show');
  if (btn.hasClass("disabled")) {
    return;
  }
  btn.append(" <i class='fa fa-spin fa-spinner'></i>");
  btn.addClass("disabled");

  $.ajax({type: 'GET', url: '/api/issues'})
    .then(function(data) {
      if (!data || !data.length) {
        return;
      }

      var rows = '';
      data.forEach(function(issue) {
        var tags = "";
        for (var tag in issue["Labels"]) {
          tags += "<span class='badge badge-default is-issue-tag' style='background-color: #" + issue["Labels"][tag] + ";'>" + tag + "</span>";
        }
        if (tags != "") {
          tags = "</br>" + tags;
        }

        rows += "<tr>" +
          "<td> <a href='" + issue["URL"] + "'>" + issue["Title"] + "</a>" + tags + "</td>" +
          "<td>" + issue["Repo"]["Owner"] + "/" + issue["Repo"]["Name"] + "</td>" +
          "<td>" + issue["Languages"].join(", ") + "</td>" +
          "</tr>";
      })
      $("#issues").append(rows);
      $('#issues_table').show();
      $('#issues_table').DataTable();
      btn.hide();
    })
    .fail(function() {
      alert('Could not get issues!');
    });
}

// checkPRs makes an API call to see how many PRs this user has completed then
// shows the results.
function checkPRs() {
  var results = $('#results');
  results.empty();

  $.ajax({type: 'GET', url: '/api/prs'})
    .then(function(data) {
      var validCount = 0;

      data.forEach(function(p, i) {
        validCount += p.Valid ? 1: 0;

        var t = $('#pr-template-' + (p.Valid ? 'valid' : 'invalid') + ' div').clone();

        t.find('.title').text(p.Title);
        t.find('.date').text(Date(p.Date));
        t.find('.repo').text(p.Repo.Owner + '/' + p.Repo.Name);

        results.prepend(t);
      });

      var message;
      if (data.length === 0) {
        // No PRs
        message = 'You have not opened any Pull Requests on public GitHub projects during October 2019.';
      } else if (validCount === 0) {
        // Some PRs but none count
        message = 'You have ' + data.length + ' Pull Request(s) but none of them are against approved repos.';
      } else if (validCount < 4) {
        // Some PRs that count but not quite 4
        message = 'Nice! You have ' + validCount + ' Pull Request(s) that count for Wichita Hacktoberfest. Keep it up!';
      } else {
        // >= 4 valid PRs! Woohoo!
        message = 'Excellent! You have hit the goal! Maybe hop in the devICT Slack and help others hit their goal too!';
      }

      results.prepend('<p class="lead">' + message + '</p>');
    })
    .fail(function() {
      alert('Could not get PRs!');
    });
}

// allowSharing is the action when a new users accepts the sponsor information
// sharing modal
function allowSharing() {
  $.ajax({url: '/api/share', type: 'put', dataType: 'json', data: 'true'})
    .then(function() {
      setShareInfoState(true);
      $('#sponsorShareModal').modal('hide');
    })
    .fail(function() {
      alert('Could not save preference. Report this as a bug!');
    });
}

// getShareInfo is loaded when the user first returns to the page
function getShareInfo() {
  $.ajax({url: '/api/share', type: 'get', dataType: 'json'})
    .then(function(share) {
      setShareInfoState(share);
    })
    .fail(function() {
      alert('Could not load preference. Report this as a bug!');
    });
}

// saveShareInfo makes an API call to update the user's preference for sharing
// information with sponsors.
function saveShareInfo(e) {
  var share = $('#share_info').prop('checked');

  setShareInfoState(undefined);
  $.ajax({url: '/api/share', type: 'put', dataType: 'json', data: JSON.stringify(share)})
    .then(function() {
      setShareInfoState(share);
    })
    .fail(function() {
      alert('Could not save preference. Report this as a bug!');
    });
}

// setShareInfoState is a helper function to set the checkbox and text display
// value of our info sharing input.
function setShareInfoState(share) {
  var state = $('#share_info_state');
  var checkbox = $('#share_info');

  switch (share) {
    case true:
      checkbox.prop('checked', true);
      state.html('Allowed.');
      state.removeClass().addClass('text-success');
      break;
    case false:
      checkbox.prop('checked', false);
      state.html('Disabled.');
      state.removeClass().addClass('text-danger');
      break;
    case undefined:
      state.html('Processing...');
      state.removeClass();
      break;
  }
  $('.share').removeClass('invisible');
}
