import moment from "moment";

export function relativeCommitTimeCalculator(commitTime) {
  let commitRelativeTime;

  const now = new Date();
  const timeDiff = moment.duration(moment(now).diff(commitTime));

  const days = Math.floor(timeDiff.asDays());
  const hours = Math.floor(timeDiff.asHours());
  const minutes = Math.floor(timeDiff.asMinutes());

  if (days > 0) {
    if (days >= 30) {
      const month = Math.floor(timeDiff.asMonths());
      commitRelativeTime =
        month === 1 ? month + " Month Ago" : month + " Months Ago";
    } else if (days >= 365) {
      const year = Math.floor(timeDiff.asYears());
      commitRelativeTime =
        year === 1 ? year + " Year Ago" : year + " Years Ago";
    } else {
      commitRelativeTime = days === 1 ? days + " Day Ago" : days + " Days Ago";
    }
  } else if (hours > 0) {
    commitRelativeTime =
      hours === 1 ? hours + " Hour Ago" : hours + " Hours Ago";
  } else if (minutes > 0) {
    commitRelativeTime =
      minutes === 1 ? minutes + " Minute Ago" : minutes + " Minutes Ago";
  } else {
    commitRelativeTime = "recent commit";
  }

  return commitRelativeTime;
}
