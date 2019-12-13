rm -rf _artifacts
mkdir _artifacts

for D in *; do
    if [ -d "${D}" ]; then
    	if [ ! "${D}" = "_artifacts" ]; then
    		zip "_artifacts/${D}.zip" "${D}"/*
		fi
    fi
done

aws s3 cp _arts s3://lalb-sample-functions --recursive --exclude "*.DS_Store*"
