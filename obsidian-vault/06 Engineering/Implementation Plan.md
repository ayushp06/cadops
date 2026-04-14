# Implementation Plan

Build thin Cobra handlers around small internal packages, keep repository watch and snapshot decision behavior outside command code, isolate pure filtering, debounce, and snapshot-selection rules for tests, and keep Git LFS lock execution separated from command formatting and validation.
