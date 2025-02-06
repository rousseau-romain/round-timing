#!/bin/bash

# Tableau d'URLs des images

urlToDl=(
"https://solomonk.fr/img/spells/373.svg"
"https://solomonk.fr/img/spells/370.svg"
"https://solomonk.fr/img/spells/364.svg"
"https://solomonk.fr/img/spells/367.svg"
"https://solomonk.fr/img/spells/366.svg"
"https://solomonk.fr/img/spells/390.svg"
"https://solomonk.fr/img/spells/391.svg"
"https://solomonk.fr/img/spells/392.svg"
"https://solomonk.fr/img/spells/393.svg"
"https://solomonk.fr/img/spells/394.svg"
"https://solomonk.fr/img/spells/395.svg"
"https://solomonk.fr/img/spells/396.svg"
"https://solomonk.fr/img/spells/397.svg"

)
idToRename=(
"130.svg"
"131.svg"
"132.svg"
"133.svg"
"134.svg"
"135.svg"
"136.svg"
"137.svg"
"138.svg"
"139.svg"
"140.svg"
"141.svg"
"142.svg"
)


# Vérifiez que les deux tableaux ont la même longueur
if [ "${#urlToDl[@]}" -ne "${#idToRename[@]}" ]; then
    echo "Error: The number of URLs does not match the number of IDs."
    exit 1
fi

# Dossier pour enregistrer les images
output_dir="./images"

# Assurez-vous que le dossier 'images' existe
mkdir -p "$output_dir"

# Fonction pour télécharger une image avec des tentatives de réessai
download_image_with_retry() {
    local url="$1"
    local filepath="$2"
    local retries=3
    local delay_time=1
    local backoff_factor=2

    for (( attempt=1; attempt<=retries; attempt++ )); do
        if curl -s -o "$filepath" "$url"; then
            echo "Downloaded $url to $filepath"
            return 0
        else
            http_status=$(curl -s -o /dev/null -w "%{http_code}" "$url")
            echo "Attempt $attempt failed with status $http_status: $url"

            if [ "$http_status" -eq 429 ]; then
                if [ "$attempt" -lt "$retries" ]; then
                    retry_after=$(( delay_time * (backoff_factor ** (attempt - 1)) ))
                    echo "Rate limited. Retrying after $retry_after seconds..."
                    sleep "$retry_after"
                else
                    echo "Failed to download $url after $retries attempts due to rate limiting."
                fi
            else
                if [ "$attempt" -lt "$retries" ]; then
                    echo "Retrying after $delay_time seconds..."
                    sleep "$delay_time"
                else
                    echo "Failed to download $url after $retries attempts."
                fi
            fi
        fi
    done

    return 1
}

# Télécharger toutes les images avec un délai entre chaque téléchargement
for i in "${!urlToDl[@]}"; do
    url="${urlToDl[$i]}"
    id="${idToRename[$i]}"
    filepath="$output_dir/$id"

    echo "Downloading $url as $id..."
    if ! download_image_with_retry "$url" "$filepath"; then
        echo "Error downloading $url"
    fi

    echo "Waiting for 1 seconds before next download..."
    sleep 1
done

echo "All downloads complete."