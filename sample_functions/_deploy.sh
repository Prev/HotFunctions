rm -rf _artifacts
mkdir _artifacts

for D in *; do
    if [ -d "${D}" ]; then
    	if [ ! "${D}" = "_artifacts" ]; then
    		zip -r "_artifacts/${D}.zip" "${D}" -x \*.pyc\*
			echo "${D}" >> _artifacts/lists.txt
		fi
    fi
done

aws s3 cp _artifacts s3://lalb-sample-functions --recursive --exclude "*.DS_Store*" --acl "public-read"
