# Nexus-Repo-Migration
This Bash script facilitates the migration of artifacts from one Sonatype Nexus repository to another. It is designed to be used in scenarios where artifacts need to be moved between Nexus instances, such as during system upgrades, migrations, or when consolidating repositories. The script automates the process of fetching artifact lists from the source Nexus, downloading each artifact, and then uploading them to the target Nexus repository.

### Key Components and Customizations:

1. **Source and Target Nexus Configuration**:
    - `SOURCE_NEXUS`, `TARGET_NEXUS`: URLs of the source and target Nexus repositories. Replace these with the actual URLs of your Nexus instances.
    - `SOURCE_REPO`, `TARGET_REPO`: Names of the source and target repositories. Change these to match the repository names in your Nexus instances.
    - `SOURCE_USER`, `SOURCE_PASSWORD`, `TARGET_USER`, `TARGET_PASSWORD`: Authentication credentials for accessing the source and target Nexus instances. Fill these in with the appropriate usernames and passwords.

2. **Initial Continuation Token**:
    - `initialContinuationToken`: An optional token used to resume artifact listing from a specific point in the source repository. This is useful for large migrations that might need to be paused and resumed. Update or clear this token based on your needs.

3. **Artifact Listing File**:
    - `ARTIFACTS_FILE`: Temporary file used to store the list of artifact paths fetched from the source Nexus. The script creates and deletes this file automatically.

### Operational Functions:

- **`fetch_artifacts`**: This function queries the source Nexus repository for a list of all artifact paths, using the Nexus REST API and handling pagination through continuation tokens. It writes the paths to the `ARTIFACTS_FILE`.

- **`migrate_artifact`**: For each artifact path listed in `ARTIFACTS_FILE`, this function downloads the artifact from the source Nexus and uploads it to the target Nexus. It handles each artifact individually, ensuring that all artifacts are migrated accurately.

### Execution Flow:

1. **Fetch Artifacts**: The script begins by calling `fetch_artifacts` to generate a comprehensive list of artifact paths in the source Nexus repository.

2. **Migrate Artifacts**: It then reads each artifact path from `ARTIFACTS_FILE` and uses `migrate_artifact` to transfer each artifact from the source to the target Nexus.

3. **Cleanup**: After all artifacts have been migrated, the script cleans up by deleting the temporary `ARTIFACTS_FILE`.

### Customization Instructions:

To use this script in your environment, ensure you replace the placeholder values for Nexus URLs, repository names, and authentication credentials with actual values relevant to your setup. Additionally, review the initial continuation token and adjust it as necessary for your migration process.
