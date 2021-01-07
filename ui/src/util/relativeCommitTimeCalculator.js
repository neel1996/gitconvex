import {
  differenceInCalendarDays,
  differenceInHours,
  differenceInMinutes,
  differenceInMonths,
  differenceInYears
} from "date-fns";

export function relativeCommitTimeCalculator(commitTime) {
  let commitRelativeTime;

  const days = differenceInCalendarDays(new Date(), new Date(commitTime));
  const hours = differenceInHours(new Date(), new Date(commitTime));
  const minutes = differenceInMinutes(new Date(), new Date(commitTime));

  if (days > 0) {
    if (days >= 30) {
      const month = differenceInMonths(new Date(), new Date(commitTime));
      commitRelativeTime =
        month === 1 ? month + " Month Ago" : month + " Months Ago";
    } else if (days >= 365) {
      const year = differenceInYears(new Date(), new Date(commitTime));
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
