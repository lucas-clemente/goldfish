import Ember from 'ember';
import moment from 'moment';

export function momentCalendar(params/*, hash*/) {
  if (params.length === 0) {
    throw new TypeError('Invalid Number of arguments, expected at least 1');
  }

  var res = moment.apply(this, params).calendar();
  return res.charAt(0).toLowerCase() + res.slice(1);
}

export default Ember.HTMLBars.makeBoundHelper(momentCalendar);
