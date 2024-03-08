# Helm Chart Versioner

This program is designed to automatically version a Helm chart based
on the amount of commits in the current repo and all checked out repos
in the current directory.

## Usage

```bash
./helmVersioner chart.yaml
```

## Example Usage

In a folder with a helm chart, a repo for your infra code and a checked
out repo for the app you are deploying, you can run the command and it
will create the a version for the chart that is based on the commit count.

```
.
.git/
my-chart/
my-app/
```
