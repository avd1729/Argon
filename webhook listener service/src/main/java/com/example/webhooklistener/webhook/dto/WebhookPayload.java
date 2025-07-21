package com.example.webhooklistener.webhook.dto;

import java.io.Serializable;

public class WebhookPayload implements Serializable {
    private String repositoryUrl;
    private String repoName;
    private String branch;
    private String commitId;
    private String commitMessage;
    private String pusherName;
    private String pusherEmail;
    private String timestamp;
    private String triggerSource;
    private String projectType;

    public WebhookPayload() {
    }

    public String getRepositoryUrl() {
        return repositoryUrl;
    }

    public void setRepositoryUrl(String repositoryUrl) {
        this.repositoryUrl = repositoryUrl;
    }

    public String getRepoName() {
        return repoName;
    }

    public void setRepoName(String repoName) {
        this.repoName = repoName;
    }

    public String getBranch() {
        return branch;
    }

    public void setBranch(String branch) {
        this.branch = branch;
    }

    public String getCommitId() {
        return commitId;
    }

    public void setCommitId(String commitId) {
        this.commitId = commitId;
    }

    public String getCommitMessage() {
        return commitMessage;
    }

    public void setCommitMessage(String commitMessage) {
        this.commitMessage = commitMessage;
    }

    public String getPusherName() {
        return pusherName;
    }

    public void setPusherName(String pusherName) {
        this.pusherName = pusherName;
    }

    public String getPusherEmail() {
        return pusherEmail;
    }

    public void setPusherEmail(String pusherEmail) {
        this.pusherEmail = pusherEmail;
    }

    public String getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(String timestamp) {
        this.timestamp = timestamp;
    }

    public String getTriggerSource() {
        return triggerSource;
    }

    public void setTriggerSource(String triggerSource) {
        this.triggerSource = triggerSource;
    }

    public String getProjectType() {
        return projectType;
    }

    public void setProjectType(String projectType) {
        this.projectType = projectType;
    }
}
