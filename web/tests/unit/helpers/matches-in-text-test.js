import {
  matchesInText
} from '../../../helpers/matches-in-text';
import { module, test } from 'qunit';

module('MatchesInTextHelper');

test('highlights first match', function(assert) {
  var result = matchesInText(['foobar foo baz', 'foo']);
  assert.equal(result.toHTML(), '<mark>foo</mark>bar foo baz');
});

test('preserves capitalization', function(assert) {
  var result = matchesInText(['Foobar foo baz', 'foo']);
  assert.equal(result.toHTML(), '<mark>Foo</mark>bar foo baz');
});

test('cuts off long texts', function(assert) {
  var result = matchesInText(['Foobar foo baz 123456789012345678901234567890123456789012345678901234567890', 'foo']);
  assert.equal(result.toHTML(), '<mark>Foo</mark>bar foo baz 12345678901234567890123456789012345');
});

test('handles newlines', function(assert) {
  var result = matchesInText(['foo\nbar', 'foo']);
  assert.equal(result.toHTML(), '<mark>foo</mark><br>bar');
});
