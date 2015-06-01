import { moduleForModel, test } from 'ember-qunit';

moduleForModel('folder', 'Unit | Model | folder', {
  needs: ['model:page']
});

test('has name', function(assert) {
  var folder = this.subject({id: '|foo|bar'});
  assert.equal('bar', folder.get('name'));
});

test('has root name', function(assert) {
  var folder = this.subject({id: '|'});
  assert.equal('/', folder.get('name'));
});

test('has path', function(assert) {
  var folder = this.subject({id: '|foo|bar'});
  assert.equal('/foo/bar', folder.get('path'));
});

test('has root path', function(assert) {
  var folder = this.subject({id: '|'});
  assert.equal('/', folder.get('path'));
});
