import Page from 'page';
import Image from 'image';
import File from 'file';

var suffixToClass = {
  '.md': Page,
  '.jpg': Image,
  '.png': Image,
  '.svg': Image,
};

function fileFactory (path) {
  for (var suffix in suffixToClass) {
    if (path.indexOf(suffix, path.length - suffix.length) !== -1) {
      return suffixToClass[suffix];
    }
  }
  return File;
}

export default fileFactory;
