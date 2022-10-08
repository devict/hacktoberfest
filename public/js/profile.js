// On document load we wire up all of the events
$(function() {
  getShareInfo();
  $('#issues_show').click(loadIssues);
  $('#check').click(checkPRs);
  $('#share_info').change(saveShareInfo);
});

function loadIssues() {
  var btn = $('#issues_show');
  if (btn.prop('disabled')) {
    return;
  }
  btn.prepend('<svg class="animate-spin w-5 h-5 text-gray-400 mr-2 -ml-1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M9.387 17.548l.371 1.482c.133.533.692.97 1.242.97h1c.55 0 1.109-.437 1.242-.97l.371-1.482c.133-.533.675-.846 1.203-.694l1.467.42c.529.151 1.188-.114 1.462-.591l.5-.867c.274-.477.177-1.179-.219-1.562l-1.098-1.061c-.396-.383-.396-1.008.001-1.39l1.096-1.061c.396-.382.494-1.084.22-1.561l-.501-.867c-.275-.477-.933-.742-1.461-.591l-1.467.42c-.529.151-1.07-.161-1.204-.694l-.37-1.48c-.133-.532-.692-.969-1.242-.969h-1c-.55 0-1.109.437-1.242.97l-.37 1.48c-.134.533-.675.846-1.204.695l-1.467-.42c-.529-.152-1.188.114-1.462.59l-.5.867c-.274.477-.177 1.179.22 1.562l1.096 1.059c.395.383.395 1.008 0 1.391l-1.098 1.061c-.395.383-.494 1.085-.219 1.562l.501.867c.274.477.933.742 1.462.591l1.467-.42c.528-.153 1.07.16 1.203.693zm2.113-7.048c1.104 0 2 .895 2 2 0 1.104-.896 2-2 2s-2-.896-2-2c0-1.105.896-2 2-2z"/></svg>');
  btn.prop('disabled', true);

  $.ajax({type: 'GET', url: '/api/issues'})
    .then(function(data) {
      if (!data || !data.length) {
        return;
      }

      var rows = '';
      data.forEach(function(issue) {
        var tags = "";
        for (var tag in issue["Labels"]) {
          tags += '<span class="inline-flex items-center px-2 rounded-full text-xs font-medium leading-5 text-gray-900" style="background-color: #' + issue["Labels"][tag] + '">' + tag + '</span>';
        }

        rows += "<tr>" +
          '<td class="px-6 py-3 max-w-0 w-full whitespace-no-wrap text-sm leading-5 font-medium text-gray-900"><a href="' + issue["URL"] + '">' + issue["Title"] + '</a><div>' + tags + '</div></td>' +
          '<td class="px-6 py-3 max-w-0 w-full whitespace-no-wrap text-sm leading-5 font-medium text-gray-500">' + issue["Repo"]["Owner"] + "/" + issue["Repo"]["Name"] + "</td>";

        // some repos have no languages
        if (issue["Languages"] == null){
            rows += '<td class="px-6 py-3 max-w-0 w-full whitespace-no-wrap text-sm leading-5 font-medium text-gray-500">None</td></tr>';
        } else {
            rows += '<td class="px-6 py-3 max-w-0 w-full whitespace-no-wrap text-sm leading-5 font-medium text-gray-500">' + issue["Languages"].join(", ") + "</td></tr>";
        }
      })
      $("#issues").append(rows);
      $('#issues_table').show();
      $('#issues_table').DataTable({
        dom: "<'py-4 px-4 sm:flex sm:items-center sm:justify-between sm:px-6 lg:px-8'fr><'-my-2 overflow-x-auto'<'py-2 align-middle inline-block min-w-full't>><'py-4 px-4 bg-white border-t border-gray-200 sm:flex sm:items-center sm:justify-between sm:px-6 lg:px-8'<'sm:flex sm:items-center sm:justify-between text-sm space-x-3'i>p>"
      });
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
  results.empty()
  results.html('<div class="py-8 flex items-center justify-center""><svg class="animate-spin w-10 h-10 text-gray-400 mr-2 -ml-1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M13 5l.855 3.42 3.389-.971 1.501 2.6-2.535 2.449 2.535 2.451-1.5 2.6-3.39-.971-.855 3.422h-3l-.855-3.422-3.39.971-1.501-2.6 2.535-2.451-2.534-2.449 1.5-2.6 3.39.971.855-3.42h3m0-2h-3c-.918 0-1.718.625-1.939 1.516l-.354 1.412-1.4-.4c-.184-.053-.369-.078-.552-.078-.7 0-1.368.37-1.731 1l-1.5 2.6c-.459.796-.317 1.802.342 2.438l1.047 1.011-1.048 1.015c-.66.637-.802 1.643-.343 2.438l1.502 2.6c.363.631 1.031 1 1.731 1 .183 0 .368-.025.552-.076l1.399-.401.354 1.415c.222.885 1.022 1.51 1.94 1.51h3c.918 0 1.718-.625 1.939-1.516l.354-1.414 1.399.4c.184.053.369.077.552.077.7 0 1.368-.37 1.731-1l1.5-2.6c.459-.796.317-1.8-.342-2.438l-1.047-1.013 1.047-1.013c.66-.637.801-1.644.342-2.438l-1.5-2.6c-.365-.631-1.031-1-1.732-1-.184 0-.368.025-.551.076l-1.4.401-.354-1.413c-.22-.884-1.02-1.509-1.938-1.509zM11.5 10.5c1.104 0 2 .895 2 2 0 1.104-.896 2-2 2s-2-.896-2-2c0-1.105.896-2 2-2m0-1c-1.654 0-3 1.346-3 3s1.346 3 3 3 3-1.346 3-3-1.346-3-3-3z"/></svg></div>');

  $.ajax({type: 'GET', url: '/api/prs'})
    .then(function(data) {
      var requiredPullRequestsCount = $('#required-pull-requests-count').val();
      console.log("requiredPullRequestsCount: ", requiredPullRequestsCount);
      var validCount = 0;
      results.empty()
      data.forEach(function(p, i) {
        validCount += p.Valid ? 1: 0;

        var t = $('#pr-template-' + (p.Valid ? 'valid' : 'invalid') + ' > li').clone();

        t.find('.title').text(p.Title);
        t.find('.date').text(new Date(p.Date).toLocaleString());
        t.find('.repo').text(p.Repo.Owner + '/' + p.Repo.Name);

        results.prepend(t);
      });

      var message;
      if (data.length === 0) {
        // No PRs
        message = 'You have not opened any Pull Requests on public GitHub projects during October 2019.';
      } else if (validCount === 0) {
        // Some PRs but none counts
        var pullRequestsString = data.length == 1 ? 'Pull Request' : 'Pull Requests'
        message = 'You have ' + data.length + ' ' + pullRequestsString + ' but none of them are against approved repos.';
      } else if (validCount < requiredPullRequestsCount) {
        // Some PRs that count but not quite 2
        var pullRequestsString = validCount == 1 ? 'Pull Request' : 'Pull Requests'
        var countsString = validCount == 1 ? 'counts' : 'count' // grammar agreement with singular or plural subject
        message = 'Nice! You have <b>' + validCount + ' ' + pullRequestsString + '</b> that ' + countsString + ' for devICT Hacktoberfest. Keep it up!';
      } else {
        // >= 'requiredPullRequestsCount' valid PRs! Woohoo!
        message = 'Excellent! You have hit the goal! Maybe hop in the devICT Slack and help others hit their goal too!';
      }

      results.prepend('<li class="px-4 py-5 sm:py-6 sm:px-6 lg:px-8 xl:px-6">' + message + '</li>');
    })
    .fail(function() {
      alert('Could not get PRs!');
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
      state.html('Sharing is on');
      state.removeClass('text-gray-500').addClass('text-green-600');
      break;
    case false:
      checkbox.prop('checked', false);
      state.html('Sharing is off');
      state.removeClass('text-green-600').addClass('text-gray-500');
      break;
    case undefined:
      state.html('Processing...');
      state.removeClass('text-green-600').addClass('text-gray-500');
      break;
  }
  $('.share').removeClass('invisible');
}
