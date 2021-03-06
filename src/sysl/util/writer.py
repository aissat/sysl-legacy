"""Super smart code writer."""

import contextlib
import cStringIO
import inspect
import operator
import re
import textwrap


class Writer(object):
    """Convenience class for writing PlantUML output.

    Provides for output into two sections, head and body, which will be
    concatenated in the final result.
    """

    indent_level = property(operator.attrgetter('_indent'))
    column = property(operator.attrgetter('_column'))

    def __init__(self, autogen_lang=None, suppress_do_not_edit=False):
        self._indent = 0
        self._head = cStringIO.StringIO()
        self._body = cStringIO.StringIO()
        self._table = None
        self._newline = True
        self._column = 0
        self.increment = 4

        if autogen_lang is not None and not suppress_do_not_edit:
            self._autogen_warning(autogen_lang)

    def _autogen_warning(self, lang):
        if lang in ['c', 'java', 'jsonnet']:
            self.head(u'//////////////////////////////////////////')
            self.head(u'//                                      //')
            self.head(u'//  AUTOGENERATED CODE -- DO NOT EDIT!  //')
            self.head(u'//                                      //')
            self.head(u'//////////////////////////////////////////')
            self.head()
        elif lang in ['sh', 'python', 'sysl']:
            self.head(u'##########################################')
            self.head(u'##                                      ##')
            self.head(u'##  AUTOGENERATED CODE -- DO NOT EDIT!  ##')
            self.head(u'##                                      ##')
            self.head(u'##########################################')
            self.head()
        elif lang in ['html', 'xml']:
            self.head(u'<?xml version="1.0" encoding="UTF-8"?>')
            self.head(u'<!-- ================================= -->')
            self.head(u'<!-- AUTOGENERATED CODE - DO NOT EDIT! -->')
            self.head(u'<!-- ================================= -->')
            self.head()
        elif lang in ['plantuml']:
            self.head(u"''''''''''''''''''''''''''''''''''''''''''")
            self.head(u"''                                      ''")
            self.head(u"''  AUTOGENERATED CODE -- DO NOT EDIT!  ''")
            self.head(u"''                                      ''")
            self.head(u"''''''''''''''''''''''''''''''''''''''''''")
            self.head()
        else:
            assert lang is None, lang

    def __call__(self, fmt=u'', *args, **kwargs):
        """Write a format string with args to the output body section."""
        return self._write(1, self._body, fmt, *args, **kwargs)

    def __str__(self):
        """Return the complete output as a string."""
        return self._head.getvalue() + self._body.getvalue()

    def start(self, frame=0):
        """Write a PlantUML start marker."""
        return self._write(frame + 1, self._head, '@startuml')

    def end(self):
        """Write a PlantUML end marker."""
        return self('@enduml')

    @contextlib.contextmanager
    def uml(self):
        """Bracket an output operation in start and end markers."""
        self.start(1)
        try:
            yield
        finally:
            self.end()

    @contextlib.contextmanager
    def indent(self, depth=None, fmt=None, *args, **kwargs):
        if not isinstance(depth, int):
            args = (fmt,) + args
            fmt = depth
            depth = self.increment

        if fmt is not None:
            self(fmt, *args, **kwargs)

        self._indent += depth
        try:
            yield
        finally:
            self._indent -= depth

    @contextlib.contextmanager
    def table(self):
        self._table = []
        yield

        max_widths = []
        for row in self._table:
            cells = row.split('\037')
            if len(cells) > 1:
                widths = [len(cell) for cell in cells]
                max_widths = [max(a, b) for (a, b) in zip(max_widths, widths)]
                if len(max_widths) < len(widths):
                    max_widths.extend(widths[len(max_widths) - len(widths):])
        for row in self._table:
            cells = row.split('\037')
            if len(cells) > 1:
                print >>self._body, ''.join(
                    cell.ljust(w) for (cell, w) in zip(cells, max_widths)
                ).rstrip()
            else:
                print >>self._body, row

        self._table = None

    @contextlib.contextmanager
    def transaction(self):
        _indent = self._indent
        _head = self._head.tell()
        _body = self._body.tell()
        _table = None if self._table is None else self._table[:]
        _newline = self._newline
        _column = self._column
        increment = self.increment

        def rollback():
            self._indent = _indent
            self._head.truncate(_head)
            self._body.truncate(_body)
            self._table = _table
            self._newline = _newline
            self._column = _column
            self.increment = increment

        try:
            yield rollback
        except BaseException:
            rollback()
            raise

    def head(self, fmt=u'', *args, **kwargs):
        """Write a format string with args to the output header section."""
        indent = self._indent
        self._indent = 0
        try:
            return self._write(1, self._head, fmt, *args, **kwargs)
        finally:
            self._indent = indent

    def textwrap(self, text, **kwargs):
        """Word-wrap a block of text, with a prefix string before each line."""
        written = None
        width = 80 - self._indent
        lines = textwrap.wrap(text, width=width, **kwargs)
        for line in lines:
            written = self('{}', line)

        if (text == '' and written is None):
            written = self('{}', '| ')
        return written

    def _write(self, frames, out, fmt, *args, **kwargs):
        """Do the actual writing on behalf of the other methods."""
        fmt = unicode(fmt)
        written = None

        if not (args or kwargs):
            f = inspect.currentframe().f_back
            while frames:
                f = f.f_back
                frames -= 1
            s = u''.join(
                (eval(p, f.f_globals, f.f_locals) if i & 1 else
                 p.replace('{{', '{').replace('}}', '}'))
                for (i, p) in enumerate(re.split(ur'{\[(.+?)\]}', fmt)))
        else:
            s = fmt.format(*args, **kwargs)

        for line in s.split('\n'):
            if '\v' in line:
                clean = re.sub(ur'\v', u'', line)
                if 80 < len(clean):
                    indent = u'  ' * self.increment
                    lines = re.sub(ur' *\v', u'\n' + indent, line).split('\n')
                else:
                    lines = [clean]
            else:
                lines = [line]

            for line in lines:
                newline = not line.endswith('\x7f')
                line = line.rstrip('\x7f').encode('utf-8')
                if self._newline:
                    line = ' ' * self._indent * bool(line) + line

                if self._table is not None and out == self._body:
                    self._table.append(line)
                else:
                    line = line.rstrip()
                    out.write(line)
                    written = line
                    if newline:
                        out.write('\n')
                        written += '\n'

                if newline:
                    self._column = 0
                else:
                    self._column += len(line)
                self._newline = newline

        return written
