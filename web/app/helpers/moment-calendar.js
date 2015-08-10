import Ember from 'ember';
import moment from 'moment';

export default Ember.Helper.extend({
  compute: function (params) {
    if (params.length === 0) {
      throw new TypeError('Invalid Number of arguments, expected at least 1');
    }

    var res = moment.apply(this, params).calendar();
    return res.charAt(0).toLowerCase() + res.slice(1);
  }
});
