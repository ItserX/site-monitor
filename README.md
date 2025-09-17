# site-monitor

#Project Overview

The project is a distributed system for monitoring website availability, collecting metrics, and sending alerts. The system is built on a microservices architecture using modern monitoring technologies.
System Architecture

**Core Components**

    - CRUD Service - Service for managing websites in the database

    - REST API for website operations (Create, Read, Update, Delete)

    - Works with PostgreSQL

    - Checker Service - Website availability checking service

    - Performs periodic website checks

    - Stores check statuses in Redis

    - Alert Service - Alert notification service

    - Processes events and sends notifications via Telegram

Monitoring Infrastructure:

    - Prometheus - Metrics collection and storage system

    - Loki - Database for log storage

    - Grafana - Platform for metrics and log visualization

    - Promtail - Agent for collecting and sending logs to Loki

    - Pushgateway - Gateway for receiving metrics

    - Redis - Database for storing check statuses

    - PostgreSQL - Main database for storing website information
