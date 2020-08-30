"""Script to delete resources in FHIR."""
import json
import logging
import os

from google.auth.transport import requests
from google.oauth2 import service_account

BASE_URL = os.environ["BASE_URL"]
PROJECT_ID = os.environ["PROJECT_ID"]
CLOUD_REGION = os.environ["CLOUD_REGION"]
DATASET_ID = os.environ["DATASET_ID"]
FHIR_STORE_ID = os.environ["FHIR_STORE_ID"]
RESOURCE_TYPE = os.environ["RESOURCE_TYPE"]
URL = f"{BASE_URL}/projects/{PROJECT_ID}/locations/{CLOUD_REGION}"

logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)


def get_session():
    """Create an authorized requests session."""
    credentials = service_account.Credentials.from_service_account_file(
        os.environ["GOOGLE_APPLICATION_CREDENTIALS"]
    )
    scoped_credentials = credentials.with_scopes([os.environ["AUTH_URL"]])
    session = requests.AuthorizedSession(scoped_credentials)

    return session


def get_resource_id():
    """Get a FHIR resource id."""
    resource_path = (
        f"{URL}/datasets/{DATASET_ID}/fhirStores/"
        f"{FHIR_STORE_ID}/fhir/{RESOURCE_TYPE}"
    )

    session = get_session()

    headers = {"Content-Type": "application/fhir+json;charset=utf-8"}
    response = session.get(resource_path, headers=headers)
    response.raise_for_status()

    if response.status_code != 200:
        logging.info(f"Status code raised: {response.status_code}")
        return

    resource = response.json()
    resource_ids = []

    try:
        entires = resource["entry"]
        for entry in entires:
            resource_ids.append(entry["resource"]["id"])

        return (resource_ids, session)

    except KeyError:
        return (resource_ids, session)


def delete_resource():
    """
    Delete a FHIR resource.

    Regardless of whether the operation succeeds or
    fails, the server returns a 200 OK HTTP status code. To check that the
    resource was successfully deleted, search for or get the resource and
    see if it exists.
    """
    resource_ids, session = get_resource_id()
    return delete_resource_helper(resource_ids, session)


def delete_resource_helper(resource_ids, session):
    """Delete a FHIR resource helper."""
    count = len(resource_ids)
    deleted_resource = 0
    for resource_id in resource_ids:
        resource_path = (
            f"{URL}/datasets/{DATASET_ID}/fhirStores/"
            f"{FHIR_STORE_ID}/fhir/{RESOURCE_TYPE}/{resource_id}"
        )
        while count != 0:
            response = session.delete(resource_path)
            if response.status_code != 200:
                logging.info(f"Status code raised: {response.status_code}")
                logging.debug(f"{response.json()}")
                return
            deleted_resource += 1

    logging.info(f"{deleted_resource} out of {count} have been deleted.")


def get_single_resource(resource_id):
    """Get a FHIR resource.

    This function's main purpose is to confirm whether or not
    a resource has been deleted.
    """
    resource_path = "{}/datasets/{}/fhirStores/{}/fhir/{}/{}".format(
        URL, DATASET_ID, FHIR_STORE_ID, RESOURCE_TYPE, resource_id
    )

    session = get_session()

    headers = {"Content-Type": "application/fhir+json;charset=utf-8"}

    response = session.get(resource_path, headers=headers)
    response.raise_for_status()

    resource = response.json()

    print("Got {} resource:".format(resource["resourceType"]))
    print(json.dumps(resource, indent=2))

    return resource


if __name__ == "__main__":
    inp = input(
        f"Please confirm you want to delete {RESOURCE_TYPE} "
        f"from {URL}/datasets/{DATASET_ID}/fhirStores/"
        f"{FHIR_STORE_ID}/fhir/ [Y/N] \n"
    )
    yes_options = ["YES", "Yes", "yes", "Y", "y"]
    no_options = ["NO", "No", "no", "N", "n"]

    if inp in yes_options:
        delete_resource()

    elif inp in no_options:
        logging.info("No resource will be deleted")
        pass

    else:
        logging.info("Kindly choose a Yes or No")
