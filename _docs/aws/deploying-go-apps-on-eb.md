# Deploying Go apps on Elastic Beanstalk

I tried following both community resources as well as official documentation. _Community resources were far better in communicating concepts_. The
official documentation felt more like a ladder with some missing rungs; much of it is accurate, but the documented steps don't match what their own
platform does. Instead of documenting an existing resource, I am using them as reference while I work through the pain points in getting this application
deployed to Elastic Beanstalk and document my findings here.

## Issues I faced

* **Unable to start with non-working environment & configure to get it working**
    * Changes to deployments would fail and rollback
    * Changes to environment variables would fail and rollback
    * In this state, it was impossible to "fix" a broken environment from my perspective (I'm sure there are some arcane ways of doing this, but documentation does not cover it)
* **Unable to manually upload a new version of the code**
    * Deployment would fail and rollback
* **Application MUST run on port `5000` and binary MUST be named `application`; this was not covered in documentation aside from being shown as the port sample code used**
    * If you don't, the logs for the environment will not contain a `web.stdout.log` file (which is your application log output)

## Happy Path

These are the steps I came up with for deploying the application as pain-free as possible:

1. _TBD_
