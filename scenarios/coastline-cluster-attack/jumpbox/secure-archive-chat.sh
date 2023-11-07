#!/bin/bash
############################################################
# Help                                                     #
############################################################
Help()
{
   # Display Help
   echo "securely encrypts and decrypts chat archive"
   echo
   echo "Syntax: script [-e|d|h]"
   echo "options:"
   echo "e     encrypt file"
   echo "d     decrypt file"
   echo "h     Print this Help."
}

usage() { echo "Usage: $0 [-e FILE / -d FILE]" 1>&2; exit 1; }

############################################################
# Process the input options. Add options as needed.        #
############################################################
# Get the options
while getopts ":he:d:" option; do
   case $option in
      h) # display Help
         Help
         exit;;
      e)
         if [ ! -f "${OPTARG}" ]; then
            usage
         fi
         echo "Encrypting ${OPTARG}"
         s=$(echo "${OPTARG}" | cut -f 1 -d '.')
         echo "$s"
         openssl enc -pbkdf2 -in "${OPTARG}" -out "${s}".enc
         ;;
      d)
         if [ ! -f "${OPTARG}" ]; then
            usage
         fi
         echo "Decrypting ${OPTARG}"
         s=$(echo "${OPTARG}" | cut -f 1 -d '.')
         openssl enc -pbkdf2 -d -in "${OPTARG}" -out "${s}"
         ;;
      :)
         echo "Error: -${OPTARG} requires an argument"
         usage ;;
      *)
         usage
         exit 1;;
   esac
done