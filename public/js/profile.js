$(function() {
  $("#check").click(function() {
    var results = $("#results");
    results.empty();

    $.ajax({type: 'GET', url: '/api/prs'})
      .then(function(data) {
        var validCount = 0;

        data.forEach(function(p, i) {
          validCount += p.Valid ? 1: 0;

          var t = $("#pr-template-" + (p.Valid ? "valid" : "invalid") + " div").clone();

          t.find(".title").text(p.Title);
          t.find(".date").text(p.Date);
          t.find(".repo").text(p.Repo.Owner + "/" + p.Repo.Name);

          results.prepend(t);
        });

        var message;
        if (data.length === 0) {
          // No PRs
          message = "You have not opened any Pull Requests on public GitHub projects during October 2017.";
        } else if (validCount === 0) {
          // Some PRs but none count
          message = "You have " + data.length + " Pull Request(s) but none of them are against approved repos.";
        } else if (validCount < 4) {
          // Some PRs that count but not quite 4
          message = "Nice! You have " + validCount + " Pull Request(s) that we'll count for Wichita Hacktoberfest. Keep it up!";
        } else {
          // >= 4 valid PRs! Woohoo!
          message = "Excellent! You have hit the goal! Maybe hop in the devICT Slack and help others hit their goal too!";
        }

        results.prepend("<p class='lead'>" + message + "</p>");
      })
      .fail(function() {
        alert("Could not get PRs!");
      });
  });
});
