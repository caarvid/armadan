window.dayjs.extend(window.isoWeek);

function getWeekDates(nr) {
  if (nr <= 0) {
    return;
  }

  const firstWeekDay = window.dayjs().isoWeek(nr).isoWeekDay(2);
  const lastWeekDay = window.dayjs().isoWeek(nr).isoWeekDay(7);

  if (!firstWeekDay.isValid() || !lastWeekDay.isValid()) {
    return "*/* - */*";
  }

  return `${firstWeekDay.format("D/M")} - ${lastWeekDay.format("D/M")}`;
}
