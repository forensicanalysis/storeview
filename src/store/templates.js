// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

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
