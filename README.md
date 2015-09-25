# packer-provisioner-serverspec

Allows to run your serverspec in a packer provisioning phase

## Example

```json
{
    "builders" : [
        {
            "type": "amazon-ebs",
            "ssh_username": "centos", // should be sudoable
            "ssh_pty": true, // tty should be available after packer 0.8
            //...
        }
    ],

    "provisioners": [
        {
            "type": "shell",
            "script": "/path/to/your-provisoner.sh"
        },
        {
            "type": "serverspec",
            "source_path": "/path/to/your/serverspec-root"
        }
    ]
}
```

`"/path/to/your/serverspec-root"` should contain `Rakefile` with task `rake spec`,
and your own serverspec test cases.

```
/path/to/your/serverspec-root
├── Gemfile # Gemfile and lock file shouldn't be neccesary
├── Gemfile.lock
├── Rakefile
└── spec
    ├── localhost
    │   ├── cloud_init_spec.rb
    │   ├── hosts_spec.rb
    │   ├── packages_spec.rb
    │   └── users_spec.rb # and so on...
    └── spec_helper.rb

```

NOTE: serverspec task should be run in EXEC mode.
