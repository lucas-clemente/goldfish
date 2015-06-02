import { moduleForModel, test } from 'ember-qunit';

moduleForModel('page', 'Unit | Model | page', {
  needs: ['model:folder']
});

test('has name', function(assert) {
  var page = this.subject({id: '|foo.md'});
  assert.equal(page.get('name'), 'foo.md');
});

test('has extension', function(assert) {
  var page = this.subject({id: '|foo.md'});
  assert.equal(page.get('extension'), 'md');
});

test('has raw path', function(assert) {
  var page = this.subject({id: '|foo.md'});
  assert.equal(page.get('rawPath'), '/v2/raw/foo.md');
});

test('has folder path', function(assert) {
  var page = this.subject({id: '|foo|bar.md'});
  assert.equal(page.get('folderPath'), '/foo');
});

test('has root folder path', function(assert) {
  var page = this.subject({id: '|foo.md'});
  assert.equal(page.get('folderPath'), '/');
});

test('has filename as title', function(assert) {
  var page = this.subject({id: '|foo.md'});
  assert.equal(page.get('title'), 'foo.md');
});

test('takes title from markdown', function(assert) {
  var page = this.subject({id: '|foo.md', markdownSource: '#bar\n'});
  assert.equal(page.get('title'), 'bar');
});

test('compiles simple markdown', function(assert) {
  var page = this.subject({id: '|foo.md', markdownSource: 'foo _bar_'});
  assert.equal(page.get('compiledMarkdown'), '<p>foo <em>bar</em></p>\n');
});

test('removes top level heading when compiling markdown', function(assert) {
  var page = this.subject({id: '|foo.md', markdownSource: '# foo\nbar'});
  assert.equal(page.get('compiledMarkdown'), '<p>bar</p>\n');
});

test('replaces $$ equations when compiling markdown', function(assert) {
  var page = this.subject({id: '|foo.md', markdownSource: '$\\sin(x)$'});
  assert.equal(page.get('compiledMarkdown'), '<p><span><script type=\"math/tex\">\\sin(x)</script></span></p>\n');
});

test('replaces \\[ \\] equations when compiling markdown', function(assert) {
  var page = this.subject({id: '|foo.md', markdownSource: '\\[\\sin(x)\\]'});
  assert.equal(page.get('compiledMarkdown'), '<p><div><script type="math/tex;mode=display">\\sin(x)</script></div></p>\n');
});

test('replaces [[image]] links', function(assert) {
  var page = this.subject({
    id: '|foo.md',
    markdownSource: '[[image.png]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><div class="image"><img src="/v2/raw/image.png" title="" class="img-thumbnail" /></div></p>\n');
});

test('replaces [[file]] links', function(assert) {
  var page = this.subject({
    id: '|foo.md',
    markdownSource: '[[file.pdf]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><a href="/folders/|/pages/|file.pdf">file</a></p>\n');
});

test('replaces [[page]] links', function(assert) {
  var page = this.subject({
    id: '|foo.md',
    markdownSource: '[[page]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><a href="/folders/|/pages/|page.md">page</a></p>\n');
});

test('replaces [[foo|file]] links', function(assert) {
  var page = this.subject({
    id: '|foo.md',
    markdownSource: '[[foo|file.pdf]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><a href="/folders/|/pages/|file.pdf">foo</a></p>\n');
});

test('replaces [[foo|page]] links', function(assert) {
  var page = this.subject({
    id: '|foo.md',
    markdownSource: '[[foo|page]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><a href="/folders/|/pages/|page.md">foo</a></p>\n');
});

test('corrects links in deep folders', function(assert) {
  var page = this.subject({
    id: '|foo|bar.md',
    markdownSource: '[[page]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><a href="/folders/|foo/pages/|foo|page.md">page</a></p>\n');
});

test('corrects absolute links', function(assert) {
  var page = this.subject({
    id: '|foo|bar.md',
    markdownSource: '[[/page]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><a href="/folders/|/pages/|page.md">/page</a></p>\n');
});

test('corrects image links in deep folders', function(assert) {
  var page = this.subject({
    id: '|foo|bar.md',
    markdownSource: '[[image.png]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><div class="image"><img src="/v2/raw/foo/image.png" title="" class="img-thumbnail" /></div></p>\n');
});

test('corrects absolute links', function(assert) {
  var page = this.subject({
    id: '|foo|bar.md',
    markdownSource: '[[/image.png]]'
  });
  assert.equal(page.get('compiledMarkdown'), '<p><div class="image"><img src="/v2/raw/image.png" title="" class="img-thumbnail" /></div></p>\n');
});
