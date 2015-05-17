import Ember from 'ember';

export function matchesInText(params/*, hash*/) {
  var text = Ember.Handlebars.Utils.escapeExpression(params[0]);
  var term = Ember.Handlebars.Utils.escapeExpression(params[1]);

  var pos = text.toLowerCase().indexOf(term.toLowerCase());
  if (pos < 0) {
    return '';
  }
  var areaSize = 50;
  var area = text.slice(Math.max(pos-areaSize, 0), Math.min(pos+areaSize, text.length-1));
  area = area.replace(term, '<mark>' + term + '</mark>');
  area = area.replace(/\n/g, '<br>');
  return Ember.String.htmlSafe(area);
}

export default Ember.HTMLBars.makeBoundHelper(matchesInText);
