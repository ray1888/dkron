---
date: 2022-06-05
title: "dkron completion bash"
slug: dkron_completion_bash
url: /cli/dkron_completion_bash/
---
## dkron completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

```
	source <(dkron completion bash)
```

To load completions for every new session, execute once:

#### Linux:

	dkron completion bash > /etc/bash_completion.d/dkron

#### macOS:

	dkron completion bash > /usr/local/etc/bash_completion.d/dkron

You will need to start a new shell for this setup to take effect.


```
dkron completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file path
```

### SEE ALSO

* [dkron completion](/docs/cli/dkron_completion/)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 5-Jun-2022
