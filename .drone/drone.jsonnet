local drone = import 'lib/drone/drone.libsonnet';
local images = import 'lib/drone/images.libsonnet';
local triggers = import 'lib/drone/triggers.libsonnet';

local pipeline = drone.pipeline;
local step = drone.step;
local withInlineStep = drone.withInlineStep;
local withStep = drone.withStep;
local withSteps = drone.withSteps;

local runTests = {
  step: step('run tests', $.commands, image=$.image),
  commands: [
    'make test'
  ],
  image: images._images.testRunner,
};

[
  pipeline('test')
  + withStep(runTests.step)
  + triggers.pr
  + triggers.main,
]