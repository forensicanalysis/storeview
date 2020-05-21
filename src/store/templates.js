export const templates = {
  'file': {
    'icon': 'folder',
    'headers': [{
      text: 'Name',
      align: 'left',
      value: 'name',
    }, {
      text: 'Size',
      value: 'size',
      align: 'right'
    }, {
      text: 'Accessed',
      value: 'atime',
      align: 'right'
    }, {
      text: 'Modified',
      value: 'mtime',
      align: 'right'
    }, {
      text: 'Created',
      value: 'ctime',
      align: 'right'
    }, {
      text: 'Origin Path',
      value: 'origin',
    }],
  },
  'windows-registry-key': {
    'icon': 'file-cabinet',
    'headers': [{
      text: 'Key',
      align: 'left',
      value: 'key',
    }, {
      text: 'Modified',
      value: 'modified_time',
      align: 'right'
    }]
  },
  'process': {
    'icon': 'console',
    'headers': [{
      text: 'Name',
      align: 'left',
      value: 'name',
    }, {
      text: 'Command Line',
      align: 'left',
      value: 'command_line',
    }, {
      text: 'Created',
      align: 'left',
      value: 'created_time',
    }, {
      text: 'Return Code',
      align: 'left',
      value: 'return_code',
    }]
  },
  'info': { 'icon': 'info' },
  'report': { 'icon': 'file-alt' },
  'directory': {
    'icon': 'folder',
    'headers': [
      {
        text: 'Path',
        value: 'path',
        align: 'right'
      }, {
        text: 'Accessed',
        value: 'atime',
        align: 'right'
      }, {
        text: 'Modified',
        value: 'mtime',
        align: 'right'
      }, {
        text: 'Created',
        value: 'ctime',
        align: 'right'
      }
    ]
  },
};
