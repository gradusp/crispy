#!/usr/bin/env sh
systemctl disable --now crispy.service
systemctl daemon-reload
