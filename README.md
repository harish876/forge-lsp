# Forge-LSP - LSP Server for the Forge ETL Framework and a VS code client.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Tools Used](#tools-patterns)

## Overview
This project is aimed at creating a very basic Language Server to solve a very trivial problem. Every time I write a config driven application, I always toggle between my actual config file and the location in my code where I am accesing them. In doing so I usually
make a spelling error, and I just spend 10-15 mins debugging.

So the aim here was to provide code completion/intellisense for me to know what configs are typed for a job. This makes me confident in writing my application and makes me blazingly fast.

## Features
 1. LSP which provides Code completion for all your configs.
 2. Uses a custom Tree Sitter Parser created for .ini files using the tree-sitter CLI.
 3. Fast Config Parsing and Querying. ( Room for enormous improvements here )

## Tools Used
1. Custom made tree sitter parser for ini files - https://github.com/harish876/tree-sitter-ini
2. Added Go Bindings for the above library. Forked a Popular Go Bindings Library - https://github.com/harish876/go-tree-sitter

## Todos
1. Parse the python code using tree-sitter and add code completion on specific actions. ✅
2. Add Goto definition and add meta data information to each setting. ✅
3. Deafault Logfile location to be changed. ⭕
4. VS Code marketplace this bad boy once logging is resolved. ⭕

