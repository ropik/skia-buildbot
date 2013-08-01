#!/usr/bin/env python
# Copyright (c) 2013 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

""" Compile step for ASAN build. """

from default_build_step_utils import DefaultBuildStepUtils
from utils import shell_utils

import os


class AsanBuildStepUtils(DefaultBuildStepUtils):
  def Compile(self, target):
    # Run the asan_build script.
    shell_utils.Bash(['which', 'clang'])
    shell_utils.Bash(['clang', '--version'])
    os.environ['GYP_DEFINES'] = self._step.args['gyp_defines']
    print 'GYP_DEFINES="%s"' % os.environ['GYP_DEFINES']
    make_cmd = os.path.join('tools', 'asan_build')
    cmd = [make_cmd,
           target,
           'BUILDTYPE=%s' % self._step.configuration,
           ]

    cmd.extend(self._step.default_make_flags)
    cmd.extend(self._step.make_flags)
    shell_utils.Bash(cmd)